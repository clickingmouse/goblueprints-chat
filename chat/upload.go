package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userid")
	// file is a multipart.File interface also a io.Reader
	// header is a metadata obj
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		fmt.Printf("err1")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("err2")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		fmt.Printf("err3")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "Successful")
}
