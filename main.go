package main

import (
	"fmt"
	"html/template"
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
)

// Config struct
type Config struct {
	Env    string
	Server Server
}

// Server struct
type Server struct {
	Host      string `toml:"host"`
	Port      string `toml:"port"`
	StaticURL string `toml:"static-url"`
}

const cookieName = "_concat_key"

func main() {
	config := Config{Env: os.Getenv("ENV")}
	_, err := toml.DecodeFile("config/"+config.Env+".toml", &config)
	if err != nil {
		panic(err)
	}
	fmt.Println(config)

	router := gin.Default()

	if config.Env == "local" {
		router.StaticFS("/static", http.Dir("static"))
	}

	template, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(template)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/templates/index.html", gin.H{
			"staticUrl": config.Server.StaticURL,
		})
	})

	router.GET("/init", func(c *gin.Context) {
		value, err := c.Cookie(cookieName)
		if err != nil {
			value = ""
		}
		c.JSON(http.StatusOK, gin.H{
			"cookie":    value,
			"staticUrl": config.Server.StaticURL,
		})
	})

	router.DELETE("/cookie", func(c *gin.Context) {
		c.String(http.StatusOK, cookieName)
	})

	router.POST("/concat", func(c *gin.Context) {
		uid := xid.New().String()
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    cookieName,
			Value:   uid,
			Expires: time.Now().Add(10 * time.Minute),
		})

		dir := "static/outputs/" + uid
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

	router.Run(":" + config.Server.Port)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".html") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
