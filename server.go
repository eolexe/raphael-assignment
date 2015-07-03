package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocraft/web"
	"github.com/raphaeljlps/raphael-assignment/db"
)

var (
	settings struct {
		ListenAddress string `json:"listenAddress"`
		DatabaseUri   string `json:"databaseUri"`
	}

	man db.TaskManager
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "path", "config.json", "json config path")

	configFile, err := os.Open(configPath)
	if err != nil {
		//Fatal because if this fail, we can proceed, so exit!
		log.Fatalf("fail to read config file %s", err)
	}

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&settings); err != nil {
		log.Fatalf("parse config file failed %s", err)
	}

	dbmap, err := db.InitDB(settings.DatabaseUri)
	if err != nil {
		log.Fatal("db connection failed")
	}

	man = db.NewTaskManager(dbmap)

	router := web.New(Context{})
	router.Middleware((*Context).AuthorizationMiddleware)
	router.Middleware((*Context).CorsMiddleware)
	router.Middleware(web.LoggerMiddleware)
	router.Get("/todo/:id", (*Context).GetTask)
	router.Post("/todo", (*Context).SaveTask)
	router.Put("/todo/:id", (*Context).UpdateTask)
	router.Delete("/todo/:id", (*Context).DeleteTask)
	router.NotFound((*Context).NotFound)
	http.ListenAndServe(settings.ListenAddress, router)
}

//Context is the application context.
type Context struct {
}

func (c *Context) GetTask(rw web.ResponseWriter, req *web.Request) {
	id, err := idFromParam(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, jsonError("missing id param"))
		return
	}

	task, err := man.Get(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, jsonError("task not found"))
		return
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, jsonError("error building response"))
		return
	}

	fmt.Fprint(rw, string(jsonTask))
}

func (c *Context) SaveTask(rw web.ResponseWriter, req *web.Request) {
	var task db.Task

	decoder := json.NewDecoder(req.Body)
	if decoder != nil {
		if err := decoder.Decode(&task); err == nil {
			man.Create(&task)
			rw.WriteHeader(http.StatusOK)
		} else {
			fmt.Fprint(rw, jsonError("trying decode, but "+err.Error()))
		}
	}
}

func (c *Context) UpdateTask(rw web.ResponseWriter, req *web.Request) {
	var task db.Task

	decoder := json.NewDecoder(req.Body)
	if decoder != nil {
		if err := decoder.Decode(&task); err == nil {
			man.Update(&task)
			rw.WriteHeader(http.StatusOK)
		} else {
			fmt.Fprint(rw, jsonError("trying decode, but "+err.Error()))
		}
	}
}

func (c *Context) DeleteTask(rw web.ResponseWriter, req *web.Request) {
	id, err := idFromParam(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, jsonError("missing id param"))
		return
	}

	err = man.Delete(id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, jsonError("resource not found"))
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (c *Context) NotFound(rw web.ResponseWriter, req *web.Request) {
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprint(rw, jsonError("resource not found"))

}

//AuthorizationMiddleware gets the Authorization header and verify if is valid
func (c *Context) AuthorizationMiddleware(rw web.ResponseWriter, r *web.Request,
	next web.NextMiddlewareFunc) {
	authArray := r.Header["Authorization"]
	if len(authArray) > 0 {
		authorization := strings.TrimSpace(authArray[0])
		content := strings.Split(authorization, " ")
		if len(content) > 1 {
			method := content[0]
			token := content[1]

			if method == "Bearer" && token == "testkey123" {
				next(rw, r)
			} else {
				rw.WriteHeader(http.StatusForbidden)
				fmt.Fprint(rw, jsonError("invalid token"))
			}
		} else {
			rw.WriteHeader(http.StatusForbidden)
			fmt.Fprint(rw, jsonError("missing token"))
		}
	} else {
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, jsonError("missing authorization header"))
	}
}

func (c *Context) CorsMiddleware(rw web.ResponseWriter, r *web.Request,
	next web.NextMiddlewareFunc) {
	if r.Method != "OPTIONS" {
		next(rw, r)
		return
	}

	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

//Utility
func idFromParam(req *web.Request) (id int, err error) {
	idParam := req.PathParams["id"]
	id, err = strconv.Atoi(idParam)
	return
}

func jsonError(message string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", message)
}
