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
	router.Use(config.CorsMiddleware)
	router.HandleFunc("/todo/category", HandlerTodoCategory)
}

func HandlerTodoCategory(w http.ResponseWriter, r *http.Request) {
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
				categories, err := modul.GetCategory(mconn, "category")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "All Category has been found", "data", categories)
				return

			} else {
				handler.StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			handler.StatusUnauthorized(w, "You are not authorized to access this data")
			return
		}
	}
}
