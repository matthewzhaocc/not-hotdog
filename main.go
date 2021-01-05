// recognizes image to determine if they are hot dog
// http handler
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

//import the index.html template
var templates = template.Must(template.ParseFiles("index.html"))

// render the index.html template
func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, "index.html", data)
}

// the main handler
// accepts GET and POST request
// GET request does not require any kind of query param
// GET request returns the upload page
// POST request require a image that is in a field called image.jpg
// POST request expect Content-Type: multipart-form
// POST request image max size is 10MB
// POST request returns wether the uploaded image is a hot dog
// require AWS credentials that have All Access to AWS Rekognition
func IsItHotDog(w http.ResponseWriter, r *http.Request) {
	// Handles the GET request
	if r.Method == "GET" {
		display(w, "index", nil)
		return
	}
	// Extract image from the multipart form, max 10MB
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("image.jpg")
	if err != nil {
		fmt.Fprintf(w, "Error Retriving image.jpg")
		return
	}

	defer file.Close()
	// saves image to file, image.jpg
	dst, err := os.Create("image.jpg")
	defer dst.Close()
	if err != nil {
		fmt.Fprintf(w, "internal error occured")
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// read it back to a bytearray
	f, _ := os.Open("image.jpg")
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	defer f.Close()
	// contact AWS for a session
	sess := session.New(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	// attempt to rekognize the label against AWS Rekognition DetectLabel API
	svc := rekognition.New(sess)
	res, errdetect := svc.DetectLabels(&rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: content,
		},
	})
	if errdetect != nil {
		fmt.Fprintf(w, errdetect.Error())
	}
	// make the result a big string for later parsing
	var b strings.Builder
	b.Grow(1024)
	for i := 0; i < len(res.Labels); i++ {
		fmt.Fprintf(&b, *res.Labels[i].Name)
	}
	// detect if the word hotdog is in the string
	if strings.Contains(b.String(), "Hot Dog") {
		fmt.Fprintf(w, "IT IS HOTDOG")
	} else {
		fmt.Fprintf(w, "IT IS NOT HOTDOG")
	}
}

// launches the app on port 6443
func main() {
	http.HandleFunc("/", IsItHotDog)
	http.ListenAndServe(":6443", nil)
}
