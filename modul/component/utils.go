package component

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// generate unique id for document
func GenerateUID(len int) (string, error) {
	bytes := make([]byte, len)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// insert one document to collection
func InsertOneDoc(db *mongo.Database, col string, docs any) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}

	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, err
}

// get one document from collection
func GetOneDoc(db *mongo.Database, col string, filter, result any) error {
	cols := db.Collection(col)
	err := cols.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		return fmt.Errorf("error retrieving data for filter %+v: %s", filter, err)
	}

	return nil
}

// get many document from collection
func GetManyDoc(db *mongo.Database, col string, filter, result any) error {
	cols := db.Collection(col)
	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error retrieving data for filter %+v: %s", filter, err)
	}

	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), result)
	if err != nil {
		return fmt.Errorf("error retrieving data for filter %+v: %s", filter, err)
	}

	return nil
}

// check document exist in collection
func CheckDocExist(db *mongo.Database, col string, filter any) (bool, error) {
	cols := db.Collection(col)
	if cols == nil {
		return false, errors.New("collection not found")
	}

	count, err := cols.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, fmt.Errorf("error checking document for filter %+v: %s", filter, err)
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}

// update one document in collection
func UpdateOneDoc(db *mongo.Database, col string, filter, update any) (*mongo.UpdateResult, error) {
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating data for filter %+v: %s", filter, err)
	}

	return result, nil
}

// find one and update document in collection
func FindOneAndUpdateDoc(db *mongo.Database, col string, filter, update, result any) error {
	cols := db.Collection(col)
	opt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := cols.FindOneAndUpdate(context.Background(), filter, update, opt).Decode(result)
	if err != nil {
		return fmt.Errorf("error updating data for filter %+v: %s", filter, err)
	}

	return nil
}

// delete one document in collection
func DeleteOneDoc(db *mongo.Database, col string, filter any) error {
	cols := db.Collection(col)
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("error deleting data for filter %+v: %s", filter, err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no data deleted for filter %+v", filter)
	}

	return nil
}

// delete many document in collection
func DeleteManyDoc(db *mongo.Database, col string, filter any) (*mongo.DeleteResult, error) {
	cols := db.Collection(col)
	result, err := cols.DeleteMany(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting data for filter %+v: %s", filter, err)
	}

	return result, nil
}
