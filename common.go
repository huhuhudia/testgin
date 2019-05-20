package main

import (
	"path"
)

type NotePostBody struct {
	FileName string `json:"filename" binding:"required"`
	Content string `json:"content" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetFilePathViaURL( rawpath string)string{
	fullpath := path.Join(BASEDIR, rawpath)
	return fullpath
}