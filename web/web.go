/*
 __  __       _ _      _    _____
|  \/  | __ _(_) |    / \  |  ___|
| |\/| |/ _` | | |   / _ \ | |_
| |  | | (_| | | |_ / ___ \|  _|
|_|  |_|\__,_|_|_(_)_/   \_\_|

Send mails as fuck!
Author : Kunal Dawn (kunal.dawn@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/
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
