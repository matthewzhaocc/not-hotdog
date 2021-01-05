// tests the http portions of the app
package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// tests the app for recognizing a hotdog image, expect return IT IS HOTDOG
func TestImageRecognizer(t *testing.T) {
	b, w := GenerateMultipartFormData(t, "./testfiles/hotdog.jpg")
	req, err := http.NewRequest("POST", "/", &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IsItHotDog)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned a %v , expecting a http 200", status)
	}
	expected := "IT IS HOTDOG"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// tests the app for recognizing a not hotdog image, expect IT IS NOT HOTDOG
func TestNotHotDog(t *testing.T) {
	b, w := GenerateMultipartFormData(t, "./testfiles/nothotdog.jpg")
	req, err := http.NewRequest("POST", "/", &b)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if err != nil {
		return
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IsItHotDog)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned a %v , expecting a http 200", status)
	}

	expected := "IT IS NOT HOTDOG"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// helper function to generate a multipart form data payload
func GenerateMultipartFormData(t *testing.T, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile("image.jpg", file.Name()); err != nil {
		t.Errorf("Error with Creating Writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with IO.Copy: %v", err)
	}
	w.Close()
	return b, w
}

// opens the file and returns the file pointer
func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

// test the HTML serving ability of the application, only testing for return code
func TestHTMLServing(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IsItHotDog)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned a %v , expecting a http 200", status)
	}

}
