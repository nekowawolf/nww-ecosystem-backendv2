package web3_tools

import (
	"encoding/json"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestWeb3ToolsCategoriesJSONAndBSON(t *testing.T) {
	categories := []string{"Wallet", "Analytics"}
	var tool Web3Tools
	if err := json.Unmarshal([]byte(`{"name":"Example","categories":["Wallet","Analytics"],"chains":["Ethereum"]}`), &tool); err != nil {
		t.Fatalf("decode request: %v", err)
	}
	if !reflect.DeepEqual(tool.Categories, categories) {
		t.Fatalf("request categories = %v, want %v", tool.Categories, categories)
	}

	encodedJSON, err := json.Marshal(tool)
	if err != nil {
		t.Fatalf("encode response: %v", err)
	}
	var response map[string]interface{}
	if err := json.Unmarshal(encodedJSON, &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := response["categories"]; !ok {
		t.Fatal("response is missing categories")
	}
	if _, ok := response["category"]; ok {
		t.Fatal("response still contains singular category")
	}

	encodedBSON, err := bson.Marshal(tool)
	if err != nil {
		t.Fatalf("encode database document: %v", err)
	}
	var stored bson.M
	if err := bson.Unmarshal(encodedBSON, &stored); err != nil {
		t.Fatalf("decode database document: %v", err)
	}
	if _, ok := stored["categories"]; !ok {
		t.Fatal("database document is missing categories")
	}
	if _, ok := stored["category"]; ok {
		t.Fatal("database document still contains singular category")
	}
	var restored Web3Tools
	if err := bson.Unmarshal(encodedBSON, &restored); err != nil {
		t.Fatalf("restore database document: %v", err)
	}
	if !reflect.DeepEqual(restored.Categories, categories) {
		t.Fatalf("stored categories = %v, want %v", restored.Categories, categories)
	}
}