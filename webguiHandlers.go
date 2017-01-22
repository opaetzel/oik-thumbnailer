package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var sourceCurrentDir = "."
var destCurrentDir = "."

type GetCurrentDir struct {
	currentDir *string
}

type SetCurrentDir struct {
	currentDir *string
}

type ParentDir struct {
	currentDir *string
}

func (s GetCurrentDir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sourceCurrentDir, err := filepath.Abs(*(s.currentDir))
	*(s.currentDir) = sourceCurrentDir
	if err != nil {
		fmt.Printf("Could not read input directory: %s\n", sourceCurrentDir)
	}
	goFiles, err := ioutil.ReadDir(sourceCurrentDir)
	if err != nil {
		fmt.Printf("Could not read input directory: %s\n", sourceCurrentDir)
	}
	files := make([]MyFile, 0)
	for _, fileInfo := range goFiles {
		if strings.HasPrefix(fileInfo.Name(), ".") {
			continue
		}
		files = append(files, MyFile{fileInfo.Name(), fileInfo.Size(), fileInfo.IsDir()})
	}
	folder := Folder{*(s.currentDir), files}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(folder); err != nil {
		panic(err)
	}
}

func (p ParentDir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newDir, err := filepath.Abs(*(p.currentDir) + "/../")
	if err != nil {
		fmt.Printf("Could not read input directory: %s\n", sourceCurrentDir)
	}
	*(p.currentDir) = newDir
	GetCurrentDir{p.currentDir}.ServeHTTP(w, r)
}

func (s SetCurrentDir) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		notParsable(w, r, err)
		return
	}
	var file MyFile
	if err = json.Unmarshal(body, &file); err != nil {
		notParsable(w, r, err)
		return
	}
	if strings.HasPrefix(file.Name, "/") {
		if file.exists() {
			newDir := file.Name
			*(s.currentDir) = newDir
		}
	} else {
		*(s.currentDir) = *(s.currentDir) + "/" + file.Name
	}
	GetCurrentDir{s.currentDir}.ServeHTTP(w, r)
}

var Convert = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		notParsable(w, r, err)
		return
	}
	var cc ConvertCommand
	if err = json.Unmarshal(body, &cc); err != nil {
		notParsable(w, r, err)
		return
	}
	files, err := createPackageWNum(&(cc.SourceDir), &(cc.DestFile), cc.Width, cc.Height)
	if err != nil {
		internalError(w, r, err)
		return
	}
	if err = json.NewEncoder(w).Encode(files); err != nil {
		panic(err)
	}
})

func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 4194304))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	return body, nil
}

func notParsable(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(422)
	log.Println(err)
	apiErr := jsonErr{Code: 422, Message: "Error parsing input. See log for details."}
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		panic(err)
	}
}

func internalError(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	apiErr := jsonErr{Code: http.StatusInternalServerError, Message: "Internal server error. See log for details."}
	log.Println(err)
	if err := json.NewEncoder(w).Encode(apiErr); err != nil {
		panic(err)
	}
}
