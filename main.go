package main

import (
	"flag"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/huhuhudia/testgin/shorturl"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var URLSMap *shorturl.URLSMap
var BASEDIR string

func init(){
	pwd, err := os.Getwd()
	if err != nil{
		log.Fatalln(err)
	}
	BASEDIR = pwd
}

func main()  {
	var sharedFilePath string
	var domainName string
	var configPath string
	var password string
	flag.StringVar(&password, "password", "secrect", "write note must have passowrd")
	flag.StringVar(&domainName, "domain", "yellowbluewhite.top", "your domain name")
	flag.StringVar(&sharedFilePath, "sharedPath", "./files", "share file path")
	flag.StringVar(&configPath, "configPath", "./config", "your config file at here")
	flag.Parse()

	URLSMap = shorturl.New(domainName, configPath)
	URLSMap.LoadFromFile()
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		sig := <- sigs
		log.Println("process quit at ", time.Now(), " ", sig)
		URLSMap.Persist()
		os.Exit(0)
	}()


	r := gin.Default()
	r.Use(ShortUrlLookUp())
	r.Use(static.Serve("/share/", static.LocalFile(sharedFilePath, true)))
	r.Use(static.Serve("/static/", static.LocalFile("./static", false)))
	r.Use(static.Serve("/image/", static.LocalFile("./image", false)))
	r.Use(static.Serve("/notes/", static.LocalFile("./notes", true)))
	r.Use(static.Serve("/video/", static.LocalFile("./video", true)))
	r.GET("/write", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./static/html/edit.html")
	})
	r.GET("/view/:notes/:filename", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./static/html/view.html")
	})
	r.GET("/modify/:notes/:filename", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "./static/html/modify.html")
	})


	r.POST("/image", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		var fullFilePath, shortUrl string
		for _, file := range files{
			_, fullFilePath = GenerateURLAndFullFilePath(domainName, "image",file.Filename)
			dstFile ,err:= os.OpenFile(path.Join(BASEDIR, fullFilePath), os.O_CREATE|os.O_WRONLY, 0666)
			defer dstFile.Close()
			shortUrl = URLSMap.Set(fullFilePath)
			if err != nil{
				log.Fatalln(err)
			}
			srcFile, err := file.Open()
			if err != nil{
				log.Fatalln(err)
			}
			io.Copy(dstFile, srcFile)
		}
		c.JSON(http.StatusOK, gin.H{
			"shortUrl":shortUrl,
		})
	})
	r.PUT("/notes", func(c *gin.Context) {
		body := NotePostBody{}
		c.BindJSON(&body)
		log.Println(body)
		if body.Password != password{
			c.AbortWithStatusJSON(401, gin.H{
				"status":401,
				"reason":"password wrong",
			})
			return
		}
		filepath := path.Join(BASEDIR, "notes", body.FileName)

		file, err := os.OpenFile(filepath,os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil{
			c.AbortWithStatusJSON(500, gin.H{
				"status":500,
				"reason":"inernal server error",
			})
			return
		}
		defer file.Close()

		file.Write([]byte(body.Content) )
	})
	r.POST("/notes", func(c *gin.Context) {
		body := NotePostBody{}
		c.BindJSON(&body)
		log.Println(body)
		if body.Password != password{
			c.AbortWithStatusJSON(401, gin.H{
				"status":401,
				"reason":"password wrong",
			})
			return
		}
		filepath := path.Join(BASEDIR, "notes", body.FileName)
		if Exists(filepath){
			c.AbortWithStatusJSON(400, gin.H{
				"status":400,
				"reason":"filename conflic",
			})
			return
		}
		file, err := os.OpenFile(filepath,os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			c.AbortWithStatusJSON(500, gin.H{
				"status":500,
				"reason":"inernal server error",
			})
			return
		}
		defer file.Close()

		file.Write([]byte(body.Content) )
	})
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/view/notes/index")
	})
	gin.SetMode(gin.ReleaseMode)
	r.Run(":80")
}


func Exists(name string ) bool{
	if _, err := os.Stat(name); err != nil{
		if os.IsNotExist(err){
			return false
		}
	}
	return true
}

