package router

import (
	"github.com/gorilla/mux"
	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/handler"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.Use(config.CorsMiddleware)

	router.HandleFunc("/", handler.Home)
	router.HandleFunc("/register", handler.HandlerRegister)
	router.HandleFunc("/login", handler.HandlerLogin)
	router.HandleFunc("/user", handler.HandlerUser)
	router.HandleFunc("/todo", handler.HandlerTodo)
	router.HandleFunc("/todo/category", handler.HandlerTodoCategory)
	router.HandleFunc("/todo/clear", handler.HandlerTodoClear)
	router.HandleFunc("/todo/log", handler.HandlerLogTodo)
	router.HandleFunc("/user/log", handler.HandlerLogUser)

	return router
}
