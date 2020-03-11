package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
)

func main() {
	http.HandleFunc("/", upload)
	http.HandleFunc("/upload", uploadFiles)

	http.ListenAndServe("localhost:8080", nil)
}

func upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<i>Get the first x lines of the uploaded files compressed into one file.<i>
		<br>
		Be mindful that the uploaded files should have the same extension.
		<br>

		<form action="/upload" method="POST" enctype="multipart/form-data">
			<input type="file" name="files" multiple><br>
			<input type="text" name="counter"><br>
			<input type="submit">
		</form>
	`)
}

func uploadFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	files := r.MultipartForm.File["files"]
	counter, err := strconv.Atoi(r.FormValue("counter"))
	if err != nil {
		panic(err)
	}

	fileExtension := ".txt"
	finalContent := ""
	for _, file := range files {
		fileExtension = filepath.Ext(file.Filename)
		f, _ := file.Open()
		scanner := bufio.NewScanner(f)
		fileContent := ""
		i := 0
		for scanner.Scan() {
			if i == counter {
				f.Close()
				break
			}

			fileContent += scanner.Text() + "\n"
			i++
		}

		finalContent += fileContent
	}

	finalContenBytes := []byte(finalContent)
	ioutil.WriteFile("composedFile"+fileExtension, finalContenBytes, 0644)

	fmt.Fprintf(w, `Upload completed`)
}
