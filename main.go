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

var templates = template.Must(template.ParseFiles("index.html"))

func display(w http.ResponseWriter, page string, data interface{}) {
	templates.ExecuteTemplate(w, "index.html", data)
}

func IsItHotDog(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		display(w, "index", nil)
		return
	}
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("image.jpg")
	if err != nil {
		fmt.Fprintf(w, "Error Retriving image.jpg")
		return
	}

	defer file.Close()

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
	f, _ := os.Open("image.jpg")
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	defer f.Close()
	
	sess := session.New(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	
	svc := rekognition.New(sess)
	res, errdetect := svc.DetectLabels(&rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: content,
		},
	})
	if errdetect != nil {
		fmt.Fprintf(w, "something happened")
	}
	var b strings.Builder
	b.Grow(1024)
	for i := 0; i < len(res.Labels); i++ {
		fmt.Fprintf(&b, *res.Labels[i].Name)
	}
	if strings.Contains(b.String(), "Hot Dog") {
		fmt.Fprintf(w, "IT IS HOTDOG")
	} else {
		fmt.Fprintf(w, "IT IS NOT HOTDOG")
	}
}

func main() {
	http.HandleFunc("/", IsItHotDog)
	http.ListenAndServe(":6443", nil)
}