package modul_test

import (
	"fmt"
	"testing"

	"github.com/mytodolist1/todolist_be/config"
	"github.com/mytodolist1/todolist_be/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = config.MongoConnect("MONGOSTRING", "mytodolist")

func TestGetTodoByUID(t *testing.T) {
	id := "657437ffb905cf734635c9a8"

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error ObjectIDFromHex: ", err)
	}

	todos, err := modul.GetTodoFromIDUser(mconn, "todo", ID)
	if err != nil {
		fmt.Println("Error GetTodoFromIDUser: ", err)
	}

	fmt.Println("todos: ", todos)
}

func TestGetCategory(t *testing.T) {
	todos, err := modul.GetCategory(mconn, "category")
	if err != nil {
		fmt.Println("Error GetTodoFromCategory: ", err)
	}

	fmt.Println("todos: ", todos)
}