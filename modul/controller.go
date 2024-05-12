package modul

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/mytodolist1/todolist_be/model"
	comp "github.com/mytodolist1/todolist_be/modul/component"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// user
// for user
// register user
// insert one document to collection
func Register(db *mongo.Database, col string, userdata model.User) error {
	// periksa data yang diinputkan
	if userdata.Email == "" || userdata.Phonenumber == "" || userdata.Username == "" || userdata.Password == "" || userdata.ConfirmPassword == "" {
		return fmt.Errorf("data tidak lengkap")
	}

	// periksa email yang diinputkan
	err := checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		return fmt.Errorf("email tidak valid")
	}

	// periksa apakah email dan username sudah terdaftar
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Email != "" {
		return fmt.Errorf("email sudah terdaftar")
	}
	if userExists.Username != "" {
		return fmt.Errorf("username sudah terdaftar")
	}

	// periksa nomor telepon yang diinputkan
	isValid, _ := comp.ValidatePhoneNumber(userdata.Phonenumber)
	if !isValid {
		return fmt.Errorf("nomor telepon tidak valid")
	}

	// periksa password yang diinputkan
	if len(userdata.Password) < 6 {
		return fmt.Errorf("password minimal 6 karakter")
	}

	// periksa apakah password dan username mengandung spasi
	if strings.Contains(userdata.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}
	if strings.Contains(userdata.Username, " ") {
		return fmt.Errorf("username tidak boleh mengandung spasi")
	}

	// periksa konfirmasi password
	if userdata.Password != userdata.ConfirmPassword {
		return fmt.Errorf("password dan konfirmasi password tidak sama")
	}

	// hash password
	hash, err := comp.HashPassword(userdata.Password)
	if err != nil {
		return fmt.Errorf("kesalahan saat meng-hash password: %v", err)
	}

	// insert data user
	user := bson.D{
		{Key: "email", Value: userdata.Email},
		{Key: "phonenumber", Value: userdata.Phonenumber},
		{Key: "username", Value: userdata.Username},
		{Key: "password", Value: hash},
		{Key: "role", Value: "user"},
	}

	_, err = comp.InsertOneDoc(db, col, user)
	if err != nil {
		return fmt.Errorf("SignUp: %v", err)
	}

	// kirim pesan konfirmasi ke wa
	// message := `Halo ` + userdata.Username + `\n\nIni adalah pesan konfirmasi dari MyTodoList. \nUsername: ` + userdata.Username + `\nPassword: ` + userdata.Password + `\nGunakan username dan password tersebut untuk login ke aplikasi MyTodoList. \n\nTerima kasih.`
	// err = SendWhatsAppConfirmation(message, userdata.Phonenumber)
	// if err != nil {
	// 	return fmt.Errorf("SendWhatsAppConfirmation: %v", err)
	// }

	return nil
}

// login user
func LogIn(db *mongo.Database, col string, userdata model.User) (user model.User, err error) {
	// periksa data yang diinputkan
	if userdata.Username == "" || userdata.Password == "" {
		err = fmt.Errorf("data tidak lengkap")
		return user, err
	}

	// periksa apakah username sudah terdaftar
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username == "" {
		err = fmt.Errorf("username tidak ditemukan")
		return user, err
	}

	// periksa password yang diinputkan
	if !comp.CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("password salah")
		return user, err
	}

	return userExists, nil
}

// change password user
// update one document in collection
func ChangePassword(db *mongo.Database, col string, userdata model.User) (user model.User, err error) {
	// periksa data yang diinputkan
	if userdata.Password == "" || userdata.ConfirmPassword == "" {
		err = fmt.Errorf("password tidak boleh kosong")
		return user, err
	}
	if len(userdata.Password) < 6 {
		err = fmt.Errorf("password minimal 6 karakter")
		return user, err
	}
	if strings.Contains(userdata.Password, " ") {
		err = fmt.Errorf("password tidak boleh mengandung spasi")
		return user, err
	}

	// periksa konfirmasi password
	if userdata.Password != userdata.ConfirmPassword {
		err = fmt.Errorf("password dan konfirmasi password tidak sama")
		return user, err
	}

	// periksa apakah password sama dengan password sebelumnya
	userExists, err := GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		return user, err
	}
	if comp.CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("password tidak boleh sama")
		return user, err
	}

	// hash password
	hash, _ := comp.HashPassword(userdata.Password)
	userExists.Password = hash

	// update password
	filter := bson.M{"username": userdata.Username}
	update := bson.M{
		"$set": bson.M{
			"password": userExists.Password,
		},
	}

	_, err = comp.UpdateOneDoc(db, col, filter, update)
	if err != nil {
		return user, fmt.Errorf("ChangePassword: %v", err)
	}

	return user, nil
}

// update user
// find one and update document in collection
func UpdateUser(db *mongo.Database, col string, userdata model.User) (usr model.User, status bool, err error) {
	// periksa data yang diinputkan
	if userdata.Username == "" || userdata.Email == "" || userdata.Phonenumber == "" {
		err := fmt.Errorf("data tidak boleh kosong")
		return usr, false, err
	}

	// periksa apakah data yang ingin diupdate sama dengan data sebelumnya
	userExists, err := GetUserFromID(db, col, userdata.ID)
	if err != nil {
		return usr, false, err
	}
	if userdata.Username == userExists.Username && userdata.Email == userExists.Email && userdata.Phonenumber == userExists.Phonenumber {
		err = fmt.Errorf("data yang ingin diupdate tidak boleh sama")
		return usr, false, err
	}

	// periksa email yang diinputkan
	err = checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		err = fmt.Errorf("email tidak valid")
		return usr, false, err
	}

	// periksa nomor telepon yang diinputkan
	isValid, _ := comp.ValidatePhoneNumber(userdata.Phonenumber)
	if !isValid {
		err = fmt.Errorf("nomor telepon tidak valid")
		return usr, false, err
	}

	// periksa apakah username mengandung spasi
	if strings.Contains(userdata.Username, " ") {
		err = fmt.Errorf("username tidak boleh mengandung spasi")
		return usr, false, err
	}

	filter := bson.M{"_id": userdata.ID}

	var originalUser model.User
	err = comp.GetOneDoc(db, col, filter, &originalUser)
	if err != nil {
		return usr, false, err
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "email", Value: userdata.Email},
			{Key: "username", Value: userdata.Username},
		}},
	}

	err = comp.FindOneAndUpdateDoc(db, col, filter, update, &userdata)
	if err != nil {
		return usr, false, err
	}

	err = comp.Log(db, "loguser", userdata.ID, originalUser, userdata)
	if err != nil {
		return usr, false, err
	}

	return userdata, true, nil
}

// delete user
// delete one document from collection
func DeleteUser(db *mongo.Database, col, username string) (bool, error) {
	// periksa apakah username sudah terdaftar
	_, err := GetUserFromUsername(db, col, username)
	if err != nil {
		err = fmt.Errorf("username tidak ditemukan")
		return false, err
	}

	// hapus data user
	filter := bson.M{"username": username}

	err = comp.DeleteOneDoc(db, col, filter)
	if err != nil {
		return false, fmt.Errorf("DeleteUser: %v", err)
	}

	return true, nil
}

// get user by id
// get one document from collection
func GetUserFromID(db *mongo.Database, col string, id primitive.ObjectID) (user model.User, err error) {
	// mengambil data user berdasarkan id
	filter := bson.M{"_id": id}

	err = comp.GetOneDoc(db, col, filter, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// get user by username
// get one document from collection
func GetUserFromUsername(db *mongo.Database, col, username string) (user model.User, err error) {
	// mengambil data user berdasarkan username
	filter := bson.M{"username": username}

	err = comp.GetOneDoc(db, col, filter, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// for admin
// get all user
// get many document from collection
func GetUserFromRole(db *mongo.Database, col, role string) (userlist []model.User, err error) {
	// mengambil data user berdasarkan role
	filter := bson.M{"role": role}

	err = comp.GetManyDoc(db, col, filter, &userlist)
	if err != nil {
		return userlist, err
	}

	return userlist, nil
}

// todo
// for user
// insert todo
// insert one document to collection
func InsertTodo(db *mongo.Database, col string, idUser primitive.ObjectID, r *http.Request) (todo model.Todo, err error) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	deadline := r.FormValue("deadline")
	times := r.FormValue("time")
	category := r.FormValue("category")

	// periksa data yang diinputkan
	if title == "" || description == "" || deadline == "" || times == "" || category == "" {
		err = fmt.Errorf("data tidak boleh kosong")
		return todo, err
	}

	// mengambil data user berdasarkan uid
	user, err := GetUserFromID(db, "user", idUser)
	if err != nil {
		fmt.Printf("GetUserFromToken: %v\n", err)
		return todo, err
	}

	// mengubah timestamps ke format unix milli
	timestamp := time.Now().UnixMilli()

	// mengubah data yang diinputkan ke huruf kapital
	title = cases.Title(language.Indonesian).String(title)
	description = cases.Title(language.Indonesian).String(description)
	category = cases.Title(language.Indonesian).String(category)

	// menyimpan file ke github
	var fileURL string
	file, _, err := r.FormFile("file")
	if err != nil {
		fileURL = ""
	} else {
		fileURL, err = comp.SaveFileToGithub("Febriand1", "fdirga63@gmail.com", "Image", "mytodolist", r)
		if err != nil {
			fmt.Printf("SaveFileToGithub: %v\n", err)
			return model.Todo{}, err
		}
		defer file.Close()
	}

	// insert data todo
	todoData := bson.D{
		{Key: "title", Value: title},
		{Key: "description", Value: description},
		{Key: "deadline", Value: deadline},
		{Key: "time", Value: times},
		{Key: "tags", Value: bson.D{
			{Key: "category", Value: category},
		}},
		{Key: "timestamps", Value: bson.D{
			{Key: "createdat", Value: timestamp},
			{Key: "updatedat", Value: timestamp},
		}},
		{Key: "user", Value: bson.D{
			{Key: "_id", Value: user.ID.Hex()},
		}},
	}

	// menyimpan file ke database jika fileURL tidak kosong
	if fileURL != "" {
		todoData = append(todoData, bson.E{Key: "file", Value: fileURL})
	}

	_, err = comp.InsertOneDoc(db, col, todoData)
	if err != nil {
		fmt.Printf("InsertTodo: %v\n", err)
		return todo, err
	}

	// periksa apakah category sudah terdaftar
	categories, err := CheckCategory(db, "category", category)
	if err != nil {
		fmt.Printf("CheckCategory: %v\n", err)
		return todo, err
	}

	// insert category jika belum terdaftar
	if !categories {
		_, err = InsertCategory(db, "category", model.Categories{Category: category})
		if err != nil {
			fmt.Printf("InsertCategory: %v\n", err)
			return todo, err
		}
	}

	return todo, nil
}

// update todo
// find one and update document in collection
func UpdateTodo(db *mongo.Database, col string, id primitive.ObjectID, r *http.Request) (model.Todo, bool, error) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	deadline := r.FormValue("deadline")
	times := r.FormValue("time")
	category := r.FormValue("category")
	file := r.FormValue("file")

	// periksa data yang diinputkan
	if title == "" || description == "" || deadline == "" || times == "" || category == "" {
		err := fmt.Errorf("data tidak lengkap")
		return model.Todo{}, false, err
	}

	// periksa apakah data yang ingin diupdate sama dengan data sebelumnya
	todoExists, err := GetTodoFromID(db, col, id)
	if err != nil {
		return model.Todo{}, false, err
	}
	if title == todoExists.Title && description == todoExists.Description && deadline == todoExists.Deadline && times == todoExists.Time {
		err = fmt.Errorf("silahkan update data anda")
		return model.Todo{}, false, err
	}

	// menyimpan file ke github
	var fileURL string
	files, _, err := r.FormFile("file")
	if err != nil {
		fileURL = ""
	} else {
		if file != "" {
			fileURL = file
		} else {
			fileURL, err = comp.SaveFileToGithub("Febriand1", "fdirga63@gmail.com", "Image", "mytodolist", r)
			if err != nil {
				fmt.Printf("SaveFileToGithub: %v\n", err)
				return model.Todo{}, false, err
			}
			defer files.Close()
		}
	}

	filter := bson.M{"_id": id}

	var originalTodo model.Todo
	err = comp.GetOneDoc(db, col, filter, &originalTodo)
	if err != nil {
		return model.Todo{}, false, err
	}

	// mengubah timestamps ke format unix milli
	time := time.Now().UnixMilli()

	// mengubah data yang diinputkan ke huruf kapital
	title = cases.Title(language.Indonesian).String(title)
	description = cases.Title(language.Indonesian).String(description)
	category = cases.Title(language.Indonesian).String(category)

	// update data todo
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "title", Value: title},
			{Key: "description", Value: description},
			{Key: "deadline", Value: deadline},
			{Key: "time", Value: times},
			{Key: "tags", Value: bson.D{
				{Key: "category", Value: category},
			}},
			{Key: "timestamps.updatedat", Value: time},
		}},
		{Key: "$setOnInsert", Value: bson.D{
			{Key: "timestamps.createdat", Value: todoExists.TimeStamps.CreatedAt},
		}},
	}

	// menyimpan file ke database jika fileURL tidak kosong
	if fileURL != "" {
		update = append(update, bson.E{Key: "$set", Value: bson.D{{Key: "file", Value: fileURL}}})
	}

	err = comp.FindOneAndUpdateDoc(db, col, filter, update, &todoExists)
	if err != nil {
		return model.Todo{}, false, err
	}

	// simpan log
	uid := todoExists.User.ID
	err = comp.Log(db, "logtodo", uid, originalTodo, todoExists)
	if err != nil {
		return model.Todo{}, false, err
	}

	return todoExists, true, nil
}

// delete todo
// delete one document from collection
func DeleteTodo(db *mongo.Database, col string, id primitive.ObjectID) (bool, error) {
	// hapus data todo berdasarkan id
	filter := bson.M{"_id": id}

	err := comp.DeleteOneDoc(db, col, filter)
	if err != nil {
		return false, fmt.Errorf("DeleteTodo: %v", err)
	}

	return true, nil
}

// get todo by id
// get one document from collection
func GetTodoFromID(db *mongo.Database, col string, id primitive.ObjectID) (todo model.Todo, err error) {
	// mengambil data todo berdasarkan id
	filter := bson.M{"_id": id}

	err = comp.GetOneDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	// mengambil data user berdasarkan uid
	user, err := GetUserFromID(db, "user", todo.User.ID)
	if err != nil {
		return todo, fmt.Errorf("user tidak ditemukan")
	}
	dataUser := model.User{
		ID:          user.ID,
		Username:    user.Username,
		Phonenumber: user.Phonenumber,
	}
	todo.User = dataUser

	return todo, nil
}

// get todo by uid user
// get many document from collection
func GetTodoFromIDUser(db *mongo.Database, col string, idUser primitive.ObjectID) (todo []model.Todo, err error) {
	// mengambil data todo berdasarkan uid user
	filter := bson.M{"user._id": idUser.Hex()}

	err = comp.GetManyDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	// mengambil data user berdasarkan uid
	for _, s := range todo {
		user, err := GetUserFromID(db, "user", s.User.ID)
		if err != nil {
			return todo, fmt.Errorf("user tidak ditemukan")
		}
		dataUser := model.User{
			ID:          user.ID,
			Username:    user.Username,
			Phonenumber: user.Phonenumber,
		}
		s.User = dataUser

		todo = append(todo, s)
		todo = todo[1:]
	}

	return todo, nil
}

// get todo by category
// get many document from collection
func GetTodoFromCategory(db *mongo.Database, col, category string) (todo []model.Todo, err error) {
	// mengambil data todo berdasarkan category
	filter := bson.M{"tags.category": category}

	err = comp.GetManyDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

// for admin
// get all todo
// get many document from collection
func GetTodoList(db *mongo.Database, col string) (todo []model.Todo, err error) {
	// mengambil semua data todo
	filter := bson.M{}

	err = comp.GetManyDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	// mengambil data user berdasarkan uid
	for _, s := range todo {
		user, err := GetUserFromID(db, "user", s.User.ID)
		if err != nil {
			return todo, fmt.Errorf("user tidak ditemukan")
		}
		dataUser := model.User{
			ID:       user.ID,
			Username: user.Username,
		}
		s.User = dataUser
		todo = append(todo, s)
		todo = todo[1:]
	}

	return todo, nil
}

// category
// insert category
// insert one document to collection
func InsertCategory(db *mongo.Database, col string, categoryDoc model.Categories) (insertedID primitive.ObjectID, err error) {
	// mengubah data yang diinputkan ke huruf kapital
	categories := cases.Title(language.Indonesian).String(categoryDoc.Category)

	// tambhakan data category
	category := bson.M{"category": categories}

	insertedID, err = comp.InsertOneDoc(db, col, category)
	if err != nil {
		fmt.Printf("InsertCategory: %v\n", err)
		return insertedID, err
	}

	return insertedID, nil
}

// check category
// check document exist in collection
func CheckCategory(db *mongo.Database, col, category string) (bool, error) {
	// cek apakah category sudah ada
	filter := bson.M{"category": category}

	exist, err := comp.CheckDocExist(db, col, filter)
	if err != nil {
		return exist, fmt.Errorf("CheckCategory: %v", err)
	}

	return exist, nil
}

// get category
// get many document from collection
func GetCategory(db *mongo.Database, col string) (category []model.Categories, err error) {
	// mengambil semua data category
	filter := bson.M{}

	err = comp.GetManyDoc(db, col, filter, &category)
	if err != nil {
		return category, err
	}

	return category, nil
}

// todo clear
// for user
// insert todo clear
// insert one document to collection
// todo done
func TodoClear(db *mongo.Database, col string, todoID primitive.ObjectID) (bool, error) {
	// mengambil data todo berdasarkan id
	todo, err := GetTodoFromID(db, "todo", todoID)
	if err != nil {
		fmt.Println("Error GetTodoFromID in colection", col, ":", err)
		return false, err
	}

	// mengubah timestamps ke format unix milli
	time := time.Now().UnixMilli()

	// insert data todo clear
	insert := bson.D{
		{Key: "isdone", Value: true},
		{Key: "timeclear", Value: time},
		{Key: "todo", Value: bson.D{
			{Key: "_id", Value: todo.ID.Hex()},
			{Key: "title", Value: todo.Title},
			{Key: "description", Value: todo.Description},
			{Key: "deadline", Value: todo.Deadline},
			{Key: "time", Value: todo.Time},
			{Key: "file", Value: todo.File},
			{Key: "tags", Value: bson.D{
				{Key: "category", Value: todo.Tags.Category},
			}},
			{Key: "user", Value: bson.D{
				{Key: "_id", Value: todo.User.ID.Hex()},
			}},
		}},
	}

	_, err = comp.InsertOneDoc(db, col, insert)
	if err != nil {
		fmt.Println("Error InsertOneDoc in colection", col, ":", err)
		return false, err
	}

	// hapus data todo dari collection todo
	status, err := DeleteTodo(db, "todo", todo.ID)
	if err != nil {
		fmt.Println("Error DeleteTodo in colection", col, ":", err)
		return false, err
	}

	// periiksa apakah data berhasil dihapus
	if !status {
		fmt.Println("Data tidak berhasil di pindahkan")
		return false, err
	}

	return true, nil
}

func DeleteTodoClear(db *mongo.Database, col string) error {
	millisecondsAgo := int64(30 * 24 * 60 * 60 * 1000)

	filter := bson.M{
		"timeclear": bson.M{
			"$lt": time.Now().UnixMilli() - millisecondsAgo,
		},
	}

	result, err := comp.DeleteManyDoc(db, col, filter)
	if err != nil {
		log.Println("Error DeleteTodoClear in colection", col, ":", err)
		return err
	}

	fmt.Printf("Dihapus %d data todo clear\n", result.DeletedCount)

	return nil
}

// get todo clear by uid user
// get many document from collection
func GetTodoClearFromIDUser(db *mongo.Database, col string, idUser primitive.ObjectID) (todo []model.TodoClear, err error) {
	// mengambil data todo clear berdasarkan uid user
	filter := bson.M{"todo.user._id": idUser.Hex()}

	err = comp.GetManyDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

// for admin
// get todo clear
// get many document from collection
func GetTodoClear(db *mongo.Database, col string) (todo []model.TodoClear, err error) {
	// mengambil semua data todo clear
	filter := bson.M{}

	err = comp.GetManyDoc(db, col, filter, &todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

// Log Todo
// for admin
// get all log todo
// get many document from collection
func GetLogTodoList(db *mongo.Database, col string) (logTodo []model.Log, err error) {
	// mengambil semua data log todo
	filter := bson.M{}

	err = comp.GetManyDoc(db, col, filter, &logTodo)
	if err != nil {
		return logTodo, err
	}

	return logTodo, nil
}

// for user
// get log todo by uid user
// get many document from collection
func GetLogTodoFromUID(db *mongo.Database, col string, idUser primitive.ObjectID) (logTodo []model.Log, err error) {
	// mengambil data log todo berdasarkan uid user
	filter := bson.M{"uid": idUser.Hex()}

	err = comp.GetManyDoc(db, col, filter, &logTodo)
	if err != nil {
		return logTodo, err
	}

	return logTodo, nil
}

// log user
// for admin
// get all log user
// get many document from collection
func GetLogAllUser(db *mongo.Database, col string) (logUser []model.Log, err error) {
	// mengambil semua data log user
	filter := bson.M{}

	err = comp.GetManyDoc(db, col, filter, &logUser)
	if err != nil {
		return logUser, err
	}

	return logUser, nil
}

// for user
// get log user by uid user
// get one document from collection
func GetLogUserFromUID(db *mongo.Database, col string, idUser primitive.ObjectID) (logUser model.Log, err error) {
	// mengambil data log user berdasarkan uid user
	filter := bson.M{"uid": idUser.Hex()}

	err = comp.GetOneDoc(db, col, filter, &logUser)
	if err != nil {
		return logUser, err
	}

	return logUser, nil
}
