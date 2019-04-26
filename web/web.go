package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"net/http"
	"strings"
	"sync"
)

var instance *Web
var once sync.Once

type Web struct {
	box *packr.Box
}

func GetWeb() *Web {
	once.Do(func() {
		instance = &Web{}
		instance.box = packr.New("public", "./dashboard/dist/dashboard/")
	})

	return instance
}

func (web *Web) Handle(context *gin.Context) {
	path := context.Param("path")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		path = "index.html"
	}

	data, err := web.box.Find(path)
	if err == nil {
		origin := context.Request.Header.Get("origin")
		context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Write(data)
		return
	} else {
		data, err := web.box.Find("index.html")
		if err == nil {
			origin := context.Request.Header.Get("origin")
			context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			context.Writer.Write(data)
			return
		}
	}

	context.Writer.WriteHeader(http.StatusNotFound)
	return
}
