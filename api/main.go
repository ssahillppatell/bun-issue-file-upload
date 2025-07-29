package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the multipart form data
	// The argument specifies the maximum memory in bytes to use for parsing the form.
	// Larger files will be written to temporary files on disk.
	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
	if err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Retrieve the file from the form
	// "myFile" is the name attribute of the file input in the HTML form.
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file from form: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. Create a new file on the server to save the uploaded content
	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Error creating destination file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 4. Copy the uploaded file content to the new file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Access other form fields if present
	otherField := r.FormValue("name")
	fmt.Printf("Other Field Value: %s\n", otherField)

	fmt.Fprintf(w, "File '%s' uploaded successfully and other field '%s' received!", handler.Filename, otherField)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func main() {
	os.MkdirAll("./uploads", os.ModePerm)
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	http.ListenAndServe(":8080", nil)
}
