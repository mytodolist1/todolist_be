package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email           string             `bson:"email,omitempty" json:"email,omitempty"`
	Phonenumber     string             `bson:"phonenumber,omitempty" json:"phonenumber,omitempty"`
	Username        string             `bson:"username,omitempty" json:"username,omitempty"`
	Password        string             `bson:"password,omitempty" json:"password,omitempty"`
	ConfirmPassword string             `bson:"confirmpassword,omitempty" json:"confirmpassword,omitempty"`
	Role            string             `bson:"role,omitempty" json:"role,omitempty"`
}

type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title,omitempty" json:"title,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Deadline    string             `bson:"deadline,omitempty" json:"deadline,omitempty"`
	Time        string             `bson:"time,omitempty" json:"time,omitempty"`
	File        string             `bson:"file,omitempty" json:"file,omitempty"`
	Tags        Categories         `bson:"tags,omitempty" json:"tags,omitempty"`
	TimeStamps  TimeStamps         `bson:"timestamps,omitempty" json:"timestamps,omitempty"`
	User        User               `bson:"user,omitempty" json:"user,omitempty"`
}

type TimeStamps struct {
	CreatedAt int64 `bson:"createdat,omitempty" json:"createdat,omitempty"`
	UpdatedAt int64 `bson:"updatedat,omitempty" json:"updatedat,omitempty"`
}

type Categories struct {
	Category string `bson:"category,omitempty" json:"category,omitempty"`
}

type TodoClear struct {
	IsDone    bool  `bson:"isdone,omitempty" json:"isdone,omitempty"`
	TimeClear int64 `bson:"timeclear,omitempty" json:"timeclear,omitempty"`
	Todo      Todo  `bson:"todo,omitempty" json:"todo,omitempty"`
}

type Log struct {
	TimeStamp int64            `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	UID       string           `bson:"uid,omitempty" json:"uid,omitempty"`
	Action    string           `bson:"action,omitempty" json:"action,omitempty"`
	Change    []map[string]any `bson:"change,omitempty" json:"change,omitempty"`
}

type Payload struct {
	Id   primitive.ObjectID `json:"id"`
	Role string             `json:"role"`
	Exp  time.Time          `json:"exp"`
	Iat  time.Time          `json:"iat"`
	Nbf  time.Time          `json:"nbf"`
}
