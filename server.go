package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocraft/web"
)

var settings struct {
	ListenAddress string `json:"listenAddress"`
	DatabaseUri   string `json:"databaseUri"`
}

type Context struct {
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "path", "config.json", "json config path")

	// configFile, err := os.Open(configPath)
	// if err != nil {
	// 	//Fatal because if this fail, we can proceed, so exit!
	// 	log.Fatal("fail to read config file")
	// }
	//
	// decoder := json.NewDecoder(configFile)
	// if err = decoder.Decode(settings); err != nil {
	// 	log.Fatal("parse config file failed")
	// }

	router := web.New(Context{})
	router.Middleware((*Context).AuthorizationMiddleware)
	router.Middleware(web.LoggerMiddleware)
	router.Get("/", (*Context).Root)
	http.ListenAndServe("localhost:3000", router)

}

func (c *Context) Root(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello")
}

// TODO: Add Comments, remove duplicated code
func (c *Context) AuthorizationMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	auth_array := r.Header["Authorization"]
	if len(auth_array) > 0 {
		authorization := strings.TrimSpace(auth_array[0])
		content := strings.Split(authorization, " ")
		if len(content) > 1 {
			method := content[0]
			token := content[1]

			if method == "Bearer" && token == "testkey123" {
				next(rw, r)
			} else {
				rw.WriteHeader(http.StatusForbidden)
				rw.Write([]byte(`{\"error\":\"Invalid token\"}`))
			}
		} else {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte(`{\"error\":\"Missing token\"}`))
		}
	} else {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{\"error\":\"Missing auth\"}`))
	}
}
