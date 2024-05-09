package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/handler"
	"github.com/mytodolist1/todolist_be/modul"
	"github.com/mytodolist1/todolist_be/paseto"
)

func init() {
	router := mux.NewRouter()
	router.Use(config.CorsMiddleware)
	router.HandleFunc("/login", HandlerLogin).Methods("POST", "OPTIONS")
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	router.Use(config.CorsMiddleware)
	router.HandleFunc("/login", HandlerLogin).Methods("POST", "OPTIONS")
	
	if r.Method == http.MethodPost {
		err := handler.JDecoder(w, r, &datauser)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		user, err := modul.LogIn(mconn, "user", datauser)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		tokenstring, err := paseto.Encode(user.ID.Hex(), user.Role, os.Getenv("PRIVATE_KEY"))
		if err != nil {
			handler.StatusBadRequest(w, "Gagal Encode Token : "+err.Error())
			return
		}

		handler.StatusOK(w, "User "+user.Username+" has been logged in", "token", tokenstring, "data", user)
		return
	}
}
