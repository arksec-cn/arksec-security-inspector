package client

import (
	"encoding/json"
	"testing"
)

func TestESClient_IndexRequest(t *testing.T) {
	data, _ := json.Marshal(struct{ Title string }{Title: "Test"})
	EsClient.IndexRequest("1", data)
}
