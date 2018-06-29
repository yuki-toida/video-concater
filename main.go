package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

// Config struct
type Config struct {
	Env    string
	Server Server
}

// Server struct
type Server struct {
	Host       string `toml:"host"`
	Port       string `toml:"port"`
	StaticURL  string `toml:"static-url"`
	BucketName string `toml:"bucket-name"`
	CookieName string `toml:"cookie-name"`
}

func main() {
	config := Config{Env: os.Getenv("ENV")}
	_, err := toml.DecodeFile("config/"+config.Env+".toml", &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("config : %v+\n", config)

	router := gin.Default()
	router.LoadHTMLFiles("index.html")

	if config.Env == "local" {
		router.StaticFS("/static", http.Dir("static"))
	}

	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"staticUrl": config.Server.StaticURL,
		})
	})

	router.GET("/init", func(c *gin.Context) {
		value, err := c.Cookie(config.Server.CookieName)
		if err != nil {
			value = ""
		}
		c.JSON(http.StatusOK, gin.H{
			"cookie":    value,
			"staticUrl": config.Server.StaticURL,
		})
	})

	router.DELETE("/cookie", func(c *gin.Context) {
		c.String(http.StatusOK, config.Server.CookieName)
	})

	router.POST("/concat", func(c *gin.Context) {
		uid := xid.New().String()
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    config.Server.CookieName,
			Value:   uid,
			Expires: time.Now().Add(10 * time.Minute),
		})

		dir := "static/outputs/" + uid
		mp4 := dir + ".mp4"
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

		err = exec.Command("ffmpeg", "-f", "concat", "-i", inputText, "-c", "copy", mp4).Run()
		if err != nil {
			panic(err)
		}
		os.Remove(inputText)
		os.RemoveAll(dir)

		fmt.Println(config.Env)
		fmt.Println(config.Env != "local")
		if config.Env != "local" {
			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				panic(err)
			}
			defer client.Close()

			data, err := ioutil.ReadFile(mp4)
			if err != nil {
				panic(err)
			}

			w := client.Bucket(config.Server.BucketName).Object(mp4).NewWriter(ctx)
			defer w.Close()

			if _, err := w.Write(data); err != nil {
				panic(err)
			}
			if err := w.Close(); err != nil {
				panic(err)
			}
			if err := os.Remove(mp4); err != nil {
				panic(err)
			}
		}

		c.String(http.StatusOK, uid)
	})

	router.Run(":" + config.Server.Port)
}
