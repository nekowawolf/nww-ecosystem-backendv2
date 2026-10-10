package airdrop

import (
	"errors"
	"fmt"
	"time"

	"github.com/nekowawolf/airdropv2/config"

	"github.com/nekowawolf/airdropv2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const airdropsCollection = "airdrops"

func InsertAirdrop(req AirdropAdmin) (interface{}, error) {
	if req.Status == "ended" && req.EndedAt == nil {
		now := time.Now()
		req.EndedAt = &now
	}
	req.ID = primitive.NewObjectID()

	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}

	return utils.InsertDocument(airdropsCollection, req)
}

func GetAirdrops(isPaid *bool) ([]AirdropAdmin, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropsCollection)
	filter := bson.M{}

	if isPaid != nil {
		filter["is_paid"] = *isPaid
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetAirdrops Find: %v", err)
	}

	var airdrops []AirdropAdmin
	if err = cursor.All(ctx, &airdrops); err != nil {
		return nil, fmt.Errorf("GetAirdrops All: %v", err)
	}

	return airdrops, nil
}

func GetAirdropStats() (map[string]int, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropsCollection)

	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error counting total: %v", err)
	}

	activeCount, err := collection.CountDocuments(ctx, bson.M{"status": bson.M{"$ne": "ended"}})
	if err != nil {
		return nil, fmt.Errorf("error counting active: %v", err)
	}

	endedCount, err := collection.CountDocuments(ctx, bson.M{"status": "ended"})
	if err != nil {
		return nil, fmt.Errorf("error counting ended: %v", err)
	}

	return map[string]int{
		"total":  int(totalCount),
		"active": int(activeCount),
		"ended":  int(endedCount),
	}, nil
}

func GetAirdropByID(id primitive.ObjectID) (AirdropAdmin, error) {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropsCollection)
	var airdrop AirdropAdmin
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&airdrop)
	if err != nil {
		return AirdropAdmin{}, err
	}
	return airdrop, nil
}

func UpdateAirdropByID(id primitive.ObjectID, updateData AirdropAdmin) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	filter := bson.M{"_id": id}

	currentAirdrop, err := GetAirdropByID(id)
	if err != nil {
		return fmt.Errorf("UpdateAirdropByID: failed to get current airdrop: %v", err)
	}

	updateFields := bson.M{
		"name":        updateData.Name,
		"task":        updateData.Task,
		"website":     updateData.Website,
		"level":       updateData.Level,
		"status":      updateData.Status,
		"backed":      updateData.Backed,
		"funds":       updateData.Funds,
		"supply":      updateData.Supply,
		"fdv":         updateData.Fdv,
		"market_cap":  updateData.MarketCap,
		"is_vesting":  updateData.IsVesting,
		"is_paid":     updateData.IsPaid,
		"claim_url":   updateData.ClaimURL,
		"discord":     updateData.Discord,
		"twitter":     updateData.Twitter,
		"telegram":    updateData.Telegram,
		"image_url":   updateData.ImageURL,
		"description": updateData.Description,
		"guide_url":   updateData.GuideURL,
		"price":       updateData.Price,
		"usd_income":  updateData.USDIncome,
	}

	if updateData.AddedBy != nil {
		updateFields["added_by"] = updateData.AddedBy
	}

	if updateData.Status == "ended" && currentAirdrop.Status != "ended" {
		now := time.Now()
		updateFields["ended_at"] = now
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := config.Database.Collection(airdropsCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("UpdateAirdropByID: %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("no data has been changed with the specified ID")
	}

	return nil
}

func DeleteAirdropByID(id primitive.ObjectID) error {
	ctx, cancel := utils.GetDBContext()
	defer cancel()

	collection := config.Database.Collection(airdropsCollection)
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting data for ID %s in airdrops: %s", id.Hex(), err.Error())
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("data with ID %s not found in airdrops", id.Hex())
	}

	return nil
}