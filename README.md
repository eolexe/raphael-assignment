# raphael-assignment
Small golang assignment for raphael - please use this repo to commit all of your work

###Create HTTP Rest API:
1. Use Gocraft for web handler (DONE)
2. Implement simple middleware for contains header 'Authorization: Bearer testkey123' in each request. Otherwise return 403 and json struct with error (DONE)
3. Log each request include status code (DONE)
4. Implement persistence with MySQL and Gorm (https://github.com/jinzhu/gorm) (DONE)
5. Use Goose for DB migration (https://bitbucket.org/liamstask/goose) (DONE)
6. Implement save endpoint for Task object (DONE)
7. Implement update endpoint for Task object (DONE)
8. Implement get endpoint for Task object (DONE)
9. Implement delete endpoint for Task object (just update IsDeleted field)  (DONE)
10. Use CORS (reply with header Access-Control-Allow-Origin: *) (DONE)
11. Add support for OPTION HTTP method for each endpoints  (DONE)
12. Configure daemon over simple JSON config. Specify path as process flag for daemon. Required params: ListenAddress, DatabaseUri. (DONE)


###Task:
```
type Task struct {
    Id          int64
    Title       string
    Description string
    Priority    int
    CreatedAt   int64
    UpdatedAt   int64
    CompletedAt bool
    IsDeleted   bool
    IsCompleted bool
}
```

Comments: 

All db related structs/functions are in the db.go file, because, this was a single entity api, otherwise, I would have split the connection related functions in a connection.go and the domain files in their respective folders, together with the Managers and Handlers for that endpoint. Thats why all the handlers are in the main package, as well with the middleware and utilities functions. 
