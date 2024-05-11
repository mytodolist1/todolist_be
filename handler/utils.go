package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/model"
	"github.com/mytodolist1/todolist_be/paseto"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	datauser      model.User
	datatodo      model.Todo
	datatodoclear model.TodoClear
	responseData  bson.M
)

var mconn = config.MongoConnect("MONGOSTRING", "mytodolist")

func ReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func PasetoDecode(w http.ResponseWriter, r *http.Request, headerName string) (model.Payload, error) {
	token := r.Header.Get(headerName)
	if token == "" {
		return model.Payload{}, errors.New("token is empty")
	}

	payload, err := paseto.Decode(os.Getenv("PUBLIC_KEY"), token)
	if err != nil {
		return model.Payload{}, err
	}

	return payload, nil
}

func JDecoder(w http.ResponseWriter, r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return err
	}
	return nil
}

func StatusNotFound(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusNotFound,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusMethodNotAllowed(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusMethodNotAllowed,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusBadRequest(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusBadRequest,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusUnauthorized(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusUnauthorized,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusCreated(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusCreated,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusNoContent(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusNoContent,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusConflict(w http.ResponseWriter, message string) {
	responseData = bson.M{
		"status":  http.StatusConflict,
		"message": message,
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

func StatusOK(w http.ResponseWriter, message string, data ...any) {
	responseData := bson.M{
		"status":  http.StatusOK,
		"message": message,
	}

	for i := 0; i < len(data); i++ {
		if i+1 < len(data) {
			key, ok := data[i].(string)
			if !ok {
				fmt.Println("Parameter namedata harus berupa string")
				return
			}
			responseData[key] = data[i+1]
			i++
		}
	}
	fmt.Fprint(w, ReturnStruct(responseData))
}

// func StatusOK(w http.ResponseWriter, message, namedata string, data ...any) {
// 	responseData = bson.M{
// 		"status":  http.StatusOK,
// 		"message": message,
// 	}

// 	if len(data) > 0 {
// 		responseData[namedata] = data[0]
// 	}

// 	fmt.Fprint(w, ReturnStruct(responseData))
// }
