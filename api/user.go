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
		return config.CorsMiddleware(next, "GET, PUT, DELETE")
	})
	router.HandleFunc("/user", HandlerUser)
}

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	id := r.URL.Query().Get("_id")

	switch r.Method {
	case http.MethodDelete: 
		_, err := handler.PasetoDecode(w, r, "Authorization")
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		if username == "" {
			handler.StatusBadRequest(w, "Missing 'username' parameter in the URL")
			return
		}
		datauser.Username = username

		err = handler.JDecoder(w, r, &datauser)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		status, err := modul.DeleteUser(mconn, "user", username)
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			handler.StatusConflict(w, "User "+username+" cannot be deleted because it is already deleted or does not exist")
			return
		}

		handler.StatusNoContent(w, "User "+username+" has been deleted")
		return

	case http.MethodPut:
		_, err := handler.PasetoDecode(w, r, "Authorization")
		if err != nil {
			handler.StatusBadRequest(w, err.Error())
			return
		}

		if id != "" {
			if id == "" {
				handler.StatusBadRequest(w, "Missing '_id' parameter in the URL")
				return
			}
			ID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				handler.StatusBadRequest(w, "Invalid '_id' parameter in the URL")
				return
			}
			datauser.ID = ID

			err = handler.JDecoder(w, r, &datauser)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			user, _, err := modul.UpdateUser(mconn, "user", datauser)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			handler.StatusOK(w, "User "+user.Username+" has been updated")
			return

		} else if username != "" {
			if username == "" {
				handler.StatusBadRequest(w, "Missing 'username' parameter in the URL")
				return
			}
			datauser.Username = username

			err = handler.JDecoder(w, r, &datauser)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			user, err := modul.ChangePassword(mconn, "user", datauser)
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}

			handler.StatusOK(w, "User "+user.Username+" has been updated")
			return
		}

	case http.MethodGet:
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := handler.PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				handler.StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				users, err := modul.GetUserFromRole(mconn, "user", "user")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "All User has been found", "data", users)
				return

			} else {
				handler.StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			if username != "" {
				if username == "" {
					handler.StatusBadRequest(w, "Missing 'username' parameter in the URL")
					return
				}
				datauser.Username = username

				user, err := modul.GetUserFromUsername(mconn, "user", username)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "User "+user.Username+" has been found", "data", user)
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
				datauser.ID = ID

				user, err := modul.GetUserFromID(mconn, "user", ID)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "User "+user.Username+" has been found", "data", user)
				return

			} else {
				payload, err := handler.PasetoDecode(w, r, "Authorization")
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				user, err := modul.GetUserFromID(mconn, "user", payload.Id)
				if err != nil {
					handler.StatusBadRequest(w, err.Error())
					return
				}

				handler.StatusOK(w, "User "+user.Username+" has been found", "data", user)
				return
			}
		}

	default:
		handler.StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}
