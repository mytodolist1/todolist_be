package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/handler"
	"github.com/mytodolist1/todolist_be/modul"
)

func init() {
	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return config.CorsMiddleware(next, "GET")
	})
	router.HandleFunc("/todo/log", HandlerLogTodo)
}

func HandlerLogTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handler.StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := handler.PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				logs, err := modul.GetLogTodoList(mconn, "logtodo")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "All Log Todo has been found", "data", logs)
				return

			} else {
				handler.StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			payload, err := handler.PasetoDecode(w, r, "Authorization")
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			logs, err := modul.GetLogTodoFromUID(mconn, "logtodo", payload.Id)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			handler.StatusOK(w, "Log Todo has been found", "data", logs)
			return
		}
	}
}
