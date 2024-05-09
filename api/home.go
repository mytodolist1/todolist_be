package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/handler"
	"github.com/mytodolist1/todolist_be/model"
)

var (
	datauser      model.User
	datatodo      model.Todo
	datatodoclear model.TodoClear
)

var mconn = config.MongoConnect("MONGOSTRING", "mytodolist")

func init() {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return config.CorsMiddleware(next, "GET")
	})
	router.HandleFunc("/", Home)
}

func Home(w http.ResponseWriter, r *http.Request) {
	handler.StatusOK(w, "Welcome to MyTodoList API")
}
