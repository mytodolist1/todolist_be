package paseto_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mytodolist1/todolist_be/paseto"
)

// paseto
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := paseto.GenerateKey()
	fmt.Println("Private Key: ", privateKey)
	fmt.Println("Public Key: ", publicKey)

	id := "657437ffb905cf734635c9a8"
	// ID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	fmt.Println("Error ObjectIDFromHex: ", err)
	// }
	role := "user"

	hasil, err := paseto.Encode(id, role, privateKey)
	fmt.Println("hasil: ", hasil, err)
}

func TestGenerateTokenPaseto(t *testing.T) {
	privateKey, publicKey := paseto.GenerateKey()
	fmt.Println("privateKey : ", privateKey)
	fmt.Println("publicKey : ", publicKey)

	id := "657437ffb905cf734635c9a8"
	// ID, err := primitive.ObjectIDFromHex(id)
	// if err != nil {
	// 	fmt.Println("Error ObjectIDFromHex: ", err)
	// }
	role := "user"

	tokenstring, err := paseto.Encode(id, role, privateKey)
	require.NoError(t, err)
	body, err := paseto.Decode(publicKey, tokenstring)
	fmt.Println("signed : ", tokenstring)
	fmt.Println("isi : ", body)
	require.NoError(t, err)
}

func TestDecoedTokenPaseto(t *testing.T) {
	publicKey := "89d8588fd5cfce87ee5442ee720e65588e992dbd431b976f44d0c6b80d864996"
	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTA4VDAwOjIxOjEyKzA3OjAwIiwiaWF0IjoiMjAyMy0xMi0wN1QyMjoyMToxMiswNzowMCIsImlkIjoiMjBhNzE3ODJhOTM1MzlkMSIsIm5iZiI6IjIwMjMtMTItMDdUMjI6MjE6MTIrMDc6MDAiLCJyb2xlIjoidXNlciJ9YxaT9HaoixwPjOjJDxOCqTrCXB3brJ4JgDOBP51Rn2HQM3Ww9Ltyx3WVoOMLPgIrEQvuiQzGZkfyUsSBIrJACg"

	id, err := paseto.Decode(publicKey, tokenstring)
	require.NoError(t, err)
	fmt.Println("id : ", id)
}

// func TestDecodeToken(t *testing.T) {
// 	//privateKey := "461ce0e87748fd656c518b870da217dc200fc8d3b6275dda8cf14943424bf8c49e2ece1954df1ea8b151fba59cc7cbd4fb810b69716149e1c26169227bd5b6868ac78b29e58b97d4018d66ad9aed4c608028f8e188dd976fa5f61fb46b47c37365d8d07b2b8d915ec9771904b608e6ba1a91b815f9e8aece8255a660b528287e"
// 	publicKey := "3fca58bcee37564ae23005b9aefe51b93cda7327a0831f533cae57f26ae70398"
// 	//userid := "budiman"
// 	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDA3OjU2OjI4WiIsImlhdCI6IjIwMjMtMTItMDFUMDU6NTY6MjhaIiwiaWQiOiI0Yjc3NDhkOGE1NTk5MDJiIiwibmJmIjoiMjAyMy0xMi0wMVQwNTo1NjoyOFoifQmqrT7Gb8Cy2XszKQJvrsb5RCRG_t7v2AvammbA4l2N0X7rh1QTw0D7L5HBi7FcKV7S7jDhpHoQKcYX0F1mVgY"
// 	idstring := modul.DecodeGetId(publicKey, tokenstring)
// 	if idstring == "" {
// 		fmt.Println("expire token")
// 	}
// 	fmt.Println("TestDecodewithStaticKey idstring : ", idstring)
// }

// paseto
// func TestGeneratePrivateKeyPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("Private Key: ", privateKey)
// 	fmt.Println("Public Key: ", publicKey)

// 	uid := "81381f10-cd45-42e4-a72c-642f34bdd53d"
// 	hasil, err := watoken.Encode(uid, privateKey)
// 	fmt.Println("hasil: ", hasil, err)
// }

// func TestGenerateTokenPaseto(t *testing.T) {
// 	privateKey, publicKey := watoken.GenerateKey()
// 	fmt.Println("privateKey : ", privateKey)
// 	fmt.Println("publicKey : ", publicKey)
// 	userid := "81381f10-cd45-42e4-a72c-642f34bdd53d"

// 	tokenstring, err := watoken.Encode(userid, privateKey)
// 	require.NoError(t, err)
// 	body, err := watoken.Decode(publicKey, tokenstring)
// 	fmt.Println("signed : ", tokenstring)
// 	fmt.Println("isi : ", body)
// 	require.NoError(t, err)
// }

// func TestDecoedTokenPaseto(t *testing.T) {
// 	publicKey := "3fca58bcee37564ae23005b9aefe51b93cda7327a0831f533cae57f26ae70398"
// 	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDA2OjMzOjMyWiIsImlhdCI6IjIwMjMtMTItMDFUMDQ6MzM6MzJaIiwiaWQiOiJxaXFpIiwibmJmIjoiMjAyMy0xMi0wMVQwNDozMzozMloifWNdG8-O7zBRsXlT78B8T5QEH-UlvYqUWBgSa22gAIs2noox_o5QZ-gj4if8gOYurkLa2oU7T7wHWBNwOxI0sAU"

// 	uid, err := watoken.Decode(publicKey, tokenstring)
// 	require.NoError(t, err)
// 	fmt.Println("uid : ", uid)
// }

// func TestDecodeToken(t *testing.T) {
// 	//privateKey := "461ce0e87748fd656c518b870da217dc200fc8d3b6275dda8cf14943424bf8c49e2ece1954df1ea8b151fba59cc7cbd4fb810b69716149e1c26169227bd5b6868ac78b29e58b97d4018d66ad9aed4c608028f8e188dd976fa5f61fb46b47c37365d8d07b2b8d915ec9771904b608e6ba1a91b815f9e8aece8255a660b528287e"
// 	publicKey := "3fca58bcee37564ae23005b9aefe51b93cda7327a0831f533cae57f26ae70398"
// 	//userid := "awangga"
// 	tokenstring := "v4.public.eyJleHAiOiIyMDIzLTEyLTAxVDA3OjU2OjI4WiIsImlhdCI6IjIwMjMtMTItMDFUMDU6NTY6MjhaIiwiaWQiOiI0Yjc3NDhkOGE1NTk5MDJiIiwibmJmIjoiMjAyMy0xMi0wMVQwNTo1NjoyOFoifQmqrT7Gb8Cy2XszKQJvrsb5RCRG_t7v2AvammbA4l2N0X7rh1QTw0D7L5HBi7FcKV7S7jDhpHoQKcYX0F1mVgY"
// 	idstring := watoken.DecodeGetId(publicKey, tokenstring)
// 	if idstring == "" {
// 		fmt.Println("expire token")
// 	}
// 	fmt.Println("TestWaTokenDecodewithStaticKey idstring : ", idstring)
// }

// func TestGeneratePasswordHash(t *testing.T) {
// 	password := "secret"
// 	hash, _ := modul.HashPassword(password) // ignore error for the sake of simplicity

// 	fmt.Println("Password:", password)
// 	fmt.Println("Hash:    ", hash)

// 	match := modul.CheckPasswordHash(password, hash)
// 	fmt.Println("Match:   ", match)
// }

// func TestHashFunction(t *testing.T) {
// 	// mconn := SetConnection("MONGOSTRING", "mytodolist")

// 	var userdata model.User
// 	userdata.Username = "budiman"
// 	userdata.Password = "secret"

// 	filter := bson.M{"username": userdata.Username}
// 	res := atdb.GetOneDoc[model.User](mconn, "user", filter)
// 	fmt.Println("Mongo User Result: ", res)
// 	hash, _ := modul.HashPassword(userdata.Password)
// 	fmt.Println("Hash Password : ", hash)
// 	match := modul.CheckPasswordHash(userdata.Password, res.Password)
// 	fmt.Println("Match:   ", match)
// }

// func TestIsPasswordValid(t *testing.T) {
// 	// mconn := SetConnection("MONGOSTRING", "mytodolist")
// 	var userdata model.User
// 	userdata.Username = "budiman"
// 	userdata.Password = "secret"

// 	anu := modul.IsPasswordValid(mconn, "user", userdata)
// 	fmt.Println(anu)
// }
