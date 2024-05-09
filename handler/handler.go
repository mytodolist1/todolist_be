package handler

import (
	"net/http"
	"os"

	"github.com/mytodolist1/todolist_be/modul"
	"github.com/mytodolist1/todolist_be/paseto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Home(w http.ResponseWriter, r *http.Request) {
	StatusOK(w, "Welcome to MyTodoList API")
}

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		err := JDecoder(w, r, &datauser)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		err = modul.Register(mconn, "user", datauser)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		StatusCreated(w, "User "+datauser.Username+" has been created")
		return
	}
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		err := JDecoder(w, r, &datauser)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		user, err := modul.LogIn(mconn, "user", datauser)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		tokenstring, err := paseto.Encode(user.ID.Hex(), user.Role, os.Getenv("PRIVATE_KEY"))
		if err != nil {
			StatusBadRequest(w, "Gagal Encode Token : "+err.Error())
			return
		}

		StatusOK(w, "User "+user.Username+" has been logged in", "token", tokenstring, "data", user)
		return
	}
}

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	id := r.URL.Query().Get("_id")

	switch r.Method {
	case http.MethodDelete:
		_, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		if username == "" {
			StatusBadRequest(w, "Missing 'username' parameter in the URL")
			return
		}
		datauser.Username = username

		err = JDecoder(w, r, &datauser)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		status, err := modul.DeleteUser(mconn, "user", username)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			StatusConflict(w, "User "+username+" cannot be deleted because it is already deleted or does not exist")
			return
		}

		StatusNoContent(w, "User "+username+" has been deleted")
		return

	case http.MethodPut:
		_, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		if id != "" {
			if id == "" {
				StatusBadRequest(w, "Missing '_id' parameter in the URL")
				return
			}
			ID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				StatusBadRequest(w, "Invalid '_id' parameter in the URL")
				return
			}
			datauser.ID = ID

			err = JDecoder(w, r, &datauser)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			user, _, err := modul.UpdateUser(mconn, "user", datauser)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			StatusOK(w, "User "+user.Username+" has been updated")
			return

		} else if username != "" {
			if username == "" {
				StatusBadRequest(w, "Missing 'username' parameter in the URL")
				return
			}
			datauser.Username = username

			err = JDecoder(w, r, &datauser)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			user, err := modul.ChangePassword(mconn, "user", datauser)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			StatusOK(w, "User "+user.Username+" has been updated")
			return
		}

	case http.MethodGet:
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				users, err := modul.GetUserFromRole(mconn, "user", "user")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All User has been found", "data", users)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			if username != "" {
				if username == "" {
					StatusBadRequest(w, "Missing 'username' parameter in the URL")
					return
				}
				datauser.Username = username

				user, err := modul.GetUserFromUsername(mconn, "user", username)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "User "+user.Username+" has been found", "data", user)
				return

			} else if id != "" {
				if id == "" {
					StatusBadRequest(w, "Missing '_id' parameter in the URL")
					return
				}
				ID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					StatusBadRequest(w, "Invalid '_id' parameter in the URL")
					return
				}
				datauser.ID = ID

				user, err := modul.GetUserFromID(mconn, "user", ID)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "User "+user.Username+" has been found", "data", user)
				return

			} else {
				payload, err := PasetoDecode(w, r, "Authorization")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				user, err := modul.GetUserFromID(mconn, "user", payload.Id)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "User "+user.Username+" has been found", "data", user)
				return
			}
		}

	default:
		StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}

func HandlerTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("_id")
	category := r.URL.Query().Get("category")

	switch r.Method {
	case http.MethodDelete:
		_, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		if id == "" {
			StatusBadRequest(w, "Missing '_id' parameter in the URL")
			return
		}
		ID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			StatusBadRequest(w, "Invalid '_id' parameter in the URL")
			return
		}
		datatodo.ID = ID

		status, err := modul.DeleteTodo(mconn, "todo", ID)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			StatusConflict(w, "Todo cannot be deleted because it is already deleted or does not exist")
			return
		}

		StatusNoContent(w, "Todo has been deleted")
		return

	case http.MethodPut:
		_, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		if id == "" {
			StatusBadRequest(w, "Missing '_id' parameter in the URL")
			return
		}
		ID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			StatusBadRequest(w, "Invalid '_id' parameter in the URL")
			return
		}
		datatodo.ID = ID

		_, _, err = modul.UpdateTodo(mconn, "todo", ID, r)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		StatusOK(w, "Todo has been updated")
		return

	case http.MethodPost:
		payload, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		_, err = modul.InsertTodo(mconn, "todo", payload.Id, r)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		StatusCreated(w, "Todo has been created")
		return

	case http.MethodGet:
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				todos, err := modul.GetTodoList(mconn, "todo")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All Todo has been found", "data", todos)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			if category != "" {
				if category == "" {
					StatusBadRequest(w, "Missing 'category' parameter in the URL")
					return
				}
				datatodo.Tags.Category = category

				todos, err := modul.GetTodoFromCategory(mconn, "todo", category)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "Todo has been found", "data", todos)
				return

			} else if id != "" {
				if id == "" {
					StatusBadRequest(w, "Missing '_id' parameter in the URL")
					return
				}
				ID, err := primitive.ObjectIDFromHex(id)
				if err != nil {
					StatusBadRequest(w, "Invalid '_id' parameter in the URL")
					return
				}
				datatodo.ID = ID

				todos, err := modul.GetTodoFromID(mconn, "todo", ID)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "Todo has been found", "data", todos)
				return

			} else {
				payload, err := PasetoDecode(w, r, "Authorization")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}
				todos, err := modul.GetTodoFromIDUser(mconn, "todo", payload.Id)
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "Todo has been found", "data", todos)
				return
			}
		}

	default:
		StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}

func HandlerTodoCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				categories, err := modul.GetCategory(mconn, "category")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All Category has been found", "data", categories)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}
		}
	}
}

func HandlerTodoClear(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("_id")
	header := r.Header.Get("AuthorizationA")

	switch r.Method {
	case http.MethodPost:
		_, err := PasetoDecode(w, r, "Authorization")
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		if id == "" {
			StatusBadRequest(w, "Missing '_id' parameter in the URL")
			return
		}
		ID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			StatusBadRequest(w, "Invalid '_id' parameter in the URL")
			return
		}
		datatodoclear.Todo.ID = ID

		status, err := modul.TodoClear(mconn, "todoclear", ID)
		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}
		if !status {
			StatusConflict(w, "Todo cannot be cleared because it is already cleared or does not exist")
			return
		}

		StatusCreated(w, "Todo has been cleared")
		return

	case http.MethodGet:
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				todos, err := modul.GetTodoClear(mconn, "todoclear")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All Todo has been found", "data", todos)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			payload, err := PasetoDecode(w, r, "Authorization")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			todos, err := modul.GetTodoClearFromIDUser(mconn, "todoclear", payload.Id)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			StatusOK(w, "Todo has been found", "data", todos)
			return
		}

	default:
		StatusMethodNotAllowed(w, "Method not allowed")
		return
	}
}

func HandlerLogTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				logs, err := modul.GetLogTodoList(mconn, "logtodo")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All Log Todo has been found", "data", logs)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			payload, err := PasetoDecode(w, r, "Authorization")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			logs, err := modul.GetLogTodoFromUID(mconn, "logtodo", payload.Id)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			StatusOK(w, "Log Todo has been found", "data", logs)
			return
		}
	}
}

func HandlerLogUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		StatusMethodNotAllowed(w, "Method not allowed")
		return

	} else {
		header := r.Header.Get("AuthorizationA")
		if header != "" {
			payload, err := PasetoDecode(w, r, "AuthorizationA")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}
			if payload.Role == "admin" {
				logs, err := modul.GetLogAllUser(mconn, "loguser")
				if err != nil {
					StatusBadRequest(w, err.Error())
					return
				}

				StatusOK(w, "All Log User has been found", "data", logs)
				return

			} else {
				StatusUnauthorized(w, "You are not authorized to access this data")
				return
			}

		} else {
			payload, err := PasetoDecode(w, r, "Authorization")
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			logs, err := modul.GetLogUserFromUID(mconn, "loguser", payload.Id)
			if err != nil {
				StatusBadRequest(w, err.Error())
				return
			}

			StatusOK(w, "Log User has been found", "data", logs)
			return
		}
	}
}
