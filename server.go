package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

	dbmap, err := db.InitDB("zeus:omgworked@tcp(localhost:3306)/taskdb") //read from config.json
	if err != nil {
		log.Fatal("db connection failed")
	}

	man = db.NewTaskManager(dbmap)

	router := web.New(Context{})
	router.Middleware((*Context).AuthorizationMiddleware)
	router.Middleware(web.LoggerMiddleware)
	router.Post("/todo/new", (*Context).SaveTask)
	router.Get("/todo/show/:id", (*Context).GetTask)
	router.Put("/todo/edit/:id", (*Context).UpdateTask)
	router.Delete("/todo/delete/:id", (*Context).DeleteTask)
	http.ListenAndServe("localhost:3000", router)
}

//Context is the application context.
type Context struct {
}

func (c *Context) GetTask(rw web.ResponseWriter, req *web.Request) {
	// // gorm, _ := db.InitDB("zeus:omgworked@tcp(localhost:3306)/taskdb")
	// // m := dbdb.NewTaskManager(gorm)
	// // task := db.Task{Id: 0, Title: "Hello", Description: "desc", Priority: 10, CreatedAt: 1231231, UpdatedAt: 12312411, CompletedAt: 1231231, IsDeleted: false, IsCompleted: false}
	//
	// // err := m.Create(&task)
	// // log.Println(err)

	idParam := req.PathParams["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Missing id param")
		return
	}

	task, err := man.Get(id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, "Task not found")
	}

	fmt.Fprintf(rw, "Task %d - %s\n", task.Id, task.Title)

}

func (c *Context) SaveTask(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello")
}

func (c *Context) UpdateTask(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello")
}

func (c *Context) DeleteTask(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Hello")
}

//AuthorizationMiddleware gets the Authorization header and verify if is valid
func (c *Context) AuthorizationMiddleware(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
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
