package component

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Insert Log ketika user melakukan perubahan data
func Log(db *mongo.Database, col string, uid primitive.ObjectID, original, updated any) error {
	// simpan timestamps dalam bentuk milisecond
	time := time.Now().UnixMilli()

	// simpan perubahan data
	newChangeEntry := bson.D{
		{Key: "timestamp", Value: time},
		{Key: "dataold", Value: original},
		{Key: "datanew", Value: updated},
	}

	// filter berdasarkan uid
	filter := bson.M{"uid": uid.Hex()}

	// update log
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "change", Value: newChangeEntry},
		}},
	}

	result, err := UpdateOneDoc(db, col, filter, update)
	if err != nil {
		fmt.Printf("UpdateOneDoc: %v\n", err)
		return err
	}

	// jika tidak ada data yang diupdate, maka insert log baru
	if result.MatchedCount == 0 {
		logUpdate := bson.D{
			{Key: "action", Value: "update"},
			{Key: "uid", Value: uid.Hex()},
			{Key: "change", Value: []bson.D{newChangeEntry}},
		}

		_, err = InsertOneDoc(db, col, logUpdate)
		if err != nil {
			fmt.Printf("InsertOneDoc: %v\n", err)
			return err
		}

	} else {
		err = GetOneDoc(db, col, filter, &original)
		if err != nil {
			fmt.Printf("GetOneDoc: %v\n", err)
			return err
		}
	}

	fmt.Printf("Result (Before Update): %+v\n", original)
	fmt.Printf("Result (After Update): %+v\n", updated)

	return nil
}
