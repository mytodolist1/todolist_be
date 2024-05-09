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
	router.Use(func(next http.Handler) http.Handler {
		return config.CorsMiddleware(next, "GET, POST, PUT, DELETE")
	})
	router.HandleFunc("/todo", HandlerTodo)
}

func HandlerTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("_id")
	category := r.URL.Query().Get("category")

	switch r.Method {
	case http.MethodDelete:
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
		datatodo.ID = ID

		status, err := modul.DeleteTodo(mconn, "todo", ID)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			handler.StatusConflict(w, "Todo cannot be deleted because it is already deleted or does not exist")
			return
		}

		handler.StatusNoContent(w, "Todo has been deleted")
		return

	case http.MethodPut:
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
		datatodo.ID = ID

		_, _, err = modul.UpdateTodo(mconn, "todo", ID, r)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		handler.StatusOK(w, "Todo has been updated")
		return

	case http.MethodPost:
		payload, err := handler.PasetoDecode(w, r, "Authorization")
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		_, err = modul.InsertTodo(mconn, "todo", payload.Id, r)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		handler.StatusCreated(w, "Todo has been created")
		return

	case http.MethodGet:
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := handler.PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				todos, err := modul.GetTodoList(mconn, "todo")
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
			if category != "" {
				if category == "" {
					handler.StatusBadRequest(w, "Missing 'category' parameter in the URL")
					return
				}
				datatodo.Tags.Category = category

				todos, err := modul.GetTodoFromCategory(mconn, "todo", category)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "Todo has been found", "data", todos)
				return

			} else if id != "" {
				if id == "" {
					handler.StatusBadRequest(w, "Missing '_id' parameter in the URL")
					return
				}
				ID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					handler.StatusBadRequest(w, "Invalid '_id' parameter in the URL")
					return
				}
				datatodo.ID = ID

				todos, err := modul.GetTodoFromID(mconn, "todo", ID)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "Todo has been found", "data", todos)
				return

			} else {
				payload, err := handler.PasetoDecode(w, r, "Authorization")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}
				todos, err := modul.GetTodoFromIDUser(mconn, "todo", payload.Id)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "Todo has been found", "data", todos)
				return
			}
		}

	default:
		handler.StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}
