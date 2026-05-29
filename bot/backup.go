package bot

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PerformBackup() (*bytes.Buffer, string, error) {
	ctx := context.Background()
	db := config.Database

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to list collections: %v", err)
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, collName := range collections {
		coll := db.Collection(collName)
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			return nil, "", fmt.Errorf("failed to query collection %s: %v", collName, err)
		}

		var docs []bson.M
		if err = cursor.All(ctx, &docs); err != nil {
			return nil, "", fmt.Errorf("failed to decode docs in %s: %v", collName, err)
		}

		if len(docs) == 0 {
			continue
		}

		fileWriter, err := zipWriter.Create(collName + ".json")
		if err != nil {
			return nil, "", fmt.Errorf("failed to create zip entry for %s: %v", collName, err)
		}

		fileWriter.Write([]byte("[\n"))

		for i, doc := range docs {
			extJson, err := bson.MarshalExtJSON(doc, false, false)
			if err != nil {
				return nil, "", fmt.Errorf("failed to marshal doc in %s: %v", collName, err)
			}
			
			fileWriter.Write(extJson)
			if i < len(docs)-1 {
				fileWriter.Write([]byte(",\n"))
			} else {
				fileWriter.Write([]byte("\n"))
			}
		}

		fileWriter.Write([]byte("]"))
	}

	if err := zipWriter.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close zip writer: %v", err)
	}

	// Update last backup date in the system_info collection
	sysInfoColl := db.Collection("system_info")
	_, _ = sysInfoColl.UpdateOne(
		ctx,
		bson.M{"_id": "backup_info"},
		bson.M{"$set": bson.M{"last_backup_date": time.Now()}},
		options.Update().SetUpsert(true),
	)

	filename := fmt.Sprintf("backup_%s.zip", time.Now().Format("2006-01-02_15-04-05"))
	return buf, filename, nil
}

func GetLastBackupDate() string {
	ctx := context.Background()
	db := config.Database
	sysInfoColl := db.Collection("system_info")
	
	var info struct {
		LastBackupDate time.Time `bson:"last_backup_date"`
	}
	err := sysInfoColl.FindOne(ctx, bson.M{"_id": "backup_info"}).Decode(&info)
	if err != nil {
		return "Never"
	}
	return info.LastBackupDate.Format("2006-01-02 15:04:05")
}
