/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

// Gouplo is a simple & easy-to-use fileserver written in Go (golang.org).
package main

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var cfgPath = flag.String("config", "config.json", "path to config file (in JSON format)")
var cfg config

func main() {
	flag.Parse()
	cfg = loadConfig(*cfgPath)
	http.HandleFunc("/", authHandler(indexHandler, hasher(cfg.User), hasher(cfg.Pass), cfg.Realm))
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/pub/", http.StripPrefix("/pub/", http.FileServer(http.Dir(cfg.PubDir))))
	go func() {
		err := http.ListenAndServeTLS(":"+cfg.HTTPSPort, cfg.CertPem, cfg.KeyPem, nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS:", err)
		}
	}()
	err := http.ListenAndServe(":"+cfg.HTTPPort, http.RedirectHandler("https://"+cfg.Domain+":"+cfg.HTTPSPort, 301))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// hasher hashes the given string and returns the sum as a slice of bytes.
func hasher(s string) []byte {
	val := md5.Sum([]byte(s))
	return val[:]
}

// config type contains the necessary server configuration strings.
type config struct {
	HTTPPort, HTTPSPort, IndexFile, PubDir, UpDir, User, Pass, Realm,
	Domain, CertPem, KeyPem string
}

// loadConfig loads configuration values from file.
func loadConfig(path string) (c config) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(cfg.IndexFile)).Execute(w, r.Host)
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
		for i := range files {
			file, err := files[i].Open() //open file
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dst, err := os.Create(cfg.UpDir + "/" + files[i].Filename) //ensure destination is writeable
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

// authHandler wraps a handler function to provide http basic authentication.
func authHandler(handler http.HandlerFunc, username, password []byte, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		userByt, passByt := hasher(user), hasher(pass)
		if !ok || subtle.ConstantTimeCompare(userByt,
			username) != 1 || subtle.ConstantTimeCompare(passByt, password) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}
