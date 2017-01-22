package main

import "os"

type Folder struct {
	Current string
	Files   []MyFile
}

type MyFile struct {
	Name  string
	Size  int64
	IsDir bool
}

type jsonErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ConvertCommand struct {
	SourceDir string `json:"sourceDir"`
	DestFile  string `json:"destFile"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func (f MyFile) exists() bool {
	_, err := os.Stat(f.Name)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
