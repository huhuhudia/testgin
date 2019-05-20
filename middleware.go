package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path"
	"time"
)




func LoggerWithFormatter()gin.HandlerFunc{
	return func(c *gin.Context){
		t := time.Now()
		c.Next()
		latency := time.Since(t)

		log.Printf("%s - [%s]  \" %s %s %s %v %s \" \n",
			c.ClientIP(),
			time.Now(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			latency,
			c.Request.UserAgent(),
		)
	}
}


func BenchMark() gin.HandlerFunc{
	return func(c *gin.Context){

		t := time.Now()
		c.Next()
		log.Printf("[benchmark] ip:%v time used:%v\n", c.ClientIP(), time.Since(t))
	}
}



func GenerateURLAndFullFilePath(domain string,service string,filename string)(string,string ){
	ext := path.Ext(filename)
	fullFilePath :=  fmt.Sprintf("%s/%x", service,md5.Sum([]byte(filename))) +ext
	url := path.Join(domain,  fullFilePath)

	return url,fullFilePath
}

func ShortUrlLookUp()gin.HandlerFunc{
	return func(c *gin.Context) {
		log.Println("=============")
		log.Println(c.Request.URL.Path)
		log.Println(URLSMap.Get(c.Request.URL.Path))
		log.Println("=============")
		if v, ok := URLSMap.Get(c.Request.URL.Path); ok{
			log.Println("redirect.....")
			c.Redirect(http.StatusTemporaryRedirect, v)
		} else {
			fmt.Println("Next")
			c.Next()
		}
	}
}

