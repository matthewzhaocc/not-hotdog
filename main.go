package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

func IsItHotDog(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("image.jpg")
	if err != nil {
		fmt.Fprintf(w, "Error Retriving image.jpg")
		return
	}

	defer file.Close()

	dst, err := os.Create(handler.Filename)
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
	fmt.Fprintf(w, *res.Labels[0].Name)
}

func main() {
	http.HandleFunc("/", IsItHotDog)
	http.ListenAndServe(":6443", nil)
}