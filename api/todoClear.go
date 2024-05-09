package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/handler"
	"github.com/mytodolist1/todolist_be/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	router := mux.NewRouter()
	router.Use(config.CorsMiddleware)
	router.HandleFunc("/todo/clear", HandlerTodoClear)
}

func HandlerTodoClear(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("_id")
	header := r.Header.Get("AuthorizationA")

	switch r.Method {
	case http.MethodPost:
		_, err := handler.PasetoDecode(w, r, "Authorization")
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		if id == "" {
			handler.StatusBadRequest(w, "Missing '_id' parameter in the URL")
			return
		}
		ID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			handler.StatusBadRequest(w, "Invalid '_id' parameter in the URL")
			return
		}
		datatodoclear.Todo.ID = ID

		status, err := modul.TodoClear(mconn, "todoclear", ID)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			handler.StatusConflict(w, "Todo cannot be cleared because it is already cleared or does not exist")
			return
		}

		handler.StatusCreated(w, "Todo has been cleared")
		return

	case http.MethodGet:
		if header != "" {
			payload, err := handler.PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				todos, err := modul.GetTodoClear(mconn, "todoclear")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "All Todo has been found", "data", todos)
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

			todos, err := modul.GetTodoClearFromIDUser(mconn, "todoclear", payload.Id)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			handler.StatusOK(w, "Todo has been found", "data", todos)
			return
		}

	default:
		handler.StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}
