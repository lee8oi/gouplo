/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

/*
gouplo is a simple Go-based file upload form that utilizes jQuery/Ajax.
*/
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")
var homefile = flag.String("home", "home.html", "home html template file")

var homeTempl = template.Must(template.ParseFiles(*homefile))

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.MultipartForm
		files := m.File["files"]
		for i, _ := range files {
			file, err := files[i].Open() //open file
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dst, err := os.Create("upload/" + files[i].Filename) //ensure destination is writeable
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := io.Copy(dst, file); err != nil { //write the file to destination
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		io.WriteString(w, "successful")
	}
}

func init() {
	if err := os.Mkdir("upload", 0777); err == nil {
		fmt.Println("upload dir created")
	} else {
		fmt.Println(err)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/pub/", http.StripPrefix("/pub/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(*addr, nil)
}
