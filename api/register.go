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
		return config.CorsMiddleware(next, "POST")
	})
	router.HandleFunc("/register", HandlerRegister)
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		err := handler.JDecoder(w, r, &datauser)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		err = modul.Register(mconn, "user", datauser)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		handler.StatusCreated(w, "User "+datauser.Username+" has been created")
		return
	}
}
