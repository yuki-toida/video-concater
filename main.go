package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

const cookieName = "_concat_key"
const outputDir = "outputs"

func main() {
	router := gin.Default()

	router.StaticFile("/", "index.html")
	router.Static("/assets", "./assets")
	router.StaticFS("/"+outputDir, http.Dir(outputDir))

	cookie := router.Group("/cookies")
	{
		cookie.GET("/", func(c *gin.Context) {
			value, err := c.Cookie(cookieName)
			if err != nil {
				value = ""
			}
			c.String(http.StatusOK, value)
		})
		cookie.DELETE("/", func(c *gin.Context) {
			c.String(http.StatusOK, cookieName)
		})
	}

	router.POST("/concat", func(c *gin.Context) {
		uid := xid.New().String()
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    cookieName,
			Value:   uid,
			Expires: time.Now().Add(time.Minute),
		})

		dir := outputDir + "/" + uid
		os.MkdirAll(dir, 0777)

		inputText := uid + ".txt"
		files := []string{}
		form, _ := c.MultipartForm()
		for _, v := range form.File {
			file := v[0]
			src, err := file.Open()
			if err != nil {
				panic(err)
			}
			defer src.Close()

			path := dir + "/" + file.Filename
			dst, err := os.Create(path)
			if err != nil {
				panic(err)
			}
			defer dst.Close()

			io.Copy(dst, src)
			files = append(files, "file '"+path+"'")
		}

		txt, err := os.Create(inputText)
		if err != nil {
			panic(err)
		}
		defer txt.Close()
		txt.Write(([]byte)(strings.Join(files, "\n")))

		err = exec.Command("ffmpeg", "-f", "concat", "-i", inputText, "-c", "copy", dir+".mp4").Run()
		if err != nil {
			panic(err)
		}
		os.Remove(inputText)
		os.RemoveAll(dir)

		c.String(http.StatusOK, uid)
	})

	router.Run()
}
