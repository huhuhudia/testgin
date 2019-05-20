package main

import (
	"log"
	"testing"
)

func TestGenerateURLAndRegularFileName(t *testing.T) {
	url, filename := GenerateURLAndFullFilePath("yellowbluewhite.top", "image", "test.jpg")
	log.Println("url:\n",url, "\n filename:\n", filename)
}

func TestGetFilePathViaURL(t *testing.T) {
	filepath := GetFilePathViaURL("/image/0412c29576c708cf0155e8de242169b1.jpg ")
	log.Println(filepath)
}