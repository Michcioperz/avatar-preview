package main

import (
	"log"
	"fmt"
	"mime"
	"io/ioutil"
	"net/http"
	"encoding/base64"
	"path"
)

const form = `<form action="/u" method="post" enctype="multipart/form-data"><input id="file" name="file" type="file"><input value="Send" type="submit"></form>`

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, ferr := r.FormFile("file")
	if ferr != nil {
		log.Print(ferr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error during upload. Did you choose the file?")
		return
	}
	defer file.Close()
	contents, cerr := ioutil.ReadAll(file)
	if cerr != nil {
		log.Print(cerr)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error during upload")
		return
	}
	extension := path.Ext(header.Filename)
	mimetype := mime.TypeByExtension(extension)
	if len(mimetype) < 1 {
		log.Print(header.Filename)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "File type not recognized. Make sure it has a reasonable file extension")
		return
	}
	encoded := fmt.Sprintf("data:%v;base64,%v", mimetype, base64.StdEncoding.EncodeToString(contents))
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, form)
	fmt.Fprintf(w, `<div></div><div style="border-radius: 50%%"></div><style>div { border: 1px solid black; width: 400px; height: 400px; background-image: url(%v); background-size: cover; background-position: center; }</style>`, encoded)
	return
}

func anyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprint(w, form)
	return
}

func main() {
	log.Print("launching")
	http.HandleFunc("/u", uploadHandler)
	http.HandleFunc("/", anyHandler)
	http.ListenAndServe(":9007", nil)
}
