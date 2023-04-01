package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"notes/static"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var letterRunes = []rune("abcdef123456789")

func init() {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
}
func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func md5sum(s string) string {
	data := []byte(s)
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}
func getContent(idmd5 string) string {
	if ret, err := ioutil.ReadFile("temp/" + idmd5); err != nil {
		return ""
	} else {
		return string(ret)
	}
}

func main() {
	r := gin.Default()
	r.Any("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/"+RandString(6))
	})
	r.GET("/:ID", func(ctx *gin.Context) {
		id := ctx.Param("ID")
		idmd5 := md5sum(id)
		content := getContent(idmd5)
		index := string(static.Index)
		index = strings.ReplaceAll(index, "{{Title_Notes.Dog}}", id)
		index = strings.ReplaceAll(index, "{{Content_Notes.Dog}}", content)
		if ua := ctx.Request.Header.Get("User-Agent"); len(ua) < 60 {
			ctx.String(200, content)
			return
		}
		ctx.Header("content-type", "text/html;charset=utf-8")
		ctx.String(200, index)
	})
	r.Any("/css.css", func(ctx *gin.Context) {
		ctx.Header("content-type", "text/css")
		ctx.String(200, string(static.Css))
	})
	r.Any("/script.js", func(ctx *gin.Context) {
		ctx.Header("content-type", "text/text")
		ctx.String(200, string(static.Script))
	})
	r.POST("/:ID", func(ctx *gin.Context) {
		ctx.Header("content-type", "text/html;charset=utf-8")
		id := ctx.Param("ID")
		idmd5 := md5sum(id)
		content := ctx.PostForm("text")
		if content == "" {
			if err := os.Remove("temp/" + idmd5); err != nil {
				log.Println("temp/"+idmd5, " Err:", err)
			}
		} else {
			ioutil.WriteFile("temp/"+idmd5, []byte(content), 0644)
		}
		ctx.String(200, "")
	})
	host := "127.0.0.1"
	port := "5522"
	flag.StringVar(&host, "l", "127.0.0.1", "listen")
	flag.StringVar(&port, "p", "8080", "port")
	flag.Parse()
	log.Println("Listen: ", fmt.Sprintf("%s:%s", host, port))
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("Listen Webapi Error : ", err)
		panic(err)
	}
}
