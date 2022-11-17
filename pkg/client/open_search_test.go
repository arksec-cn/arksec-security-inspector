package client

import (
	"encoding/json"
	"testing"
)

func TestOpenSearchClient_Info(t *testing.T) {
	//OsClient.Info()
	data, _ := json.Marshal(struct{ Title string }{Title: "Test"})
	OsClient.Add("1", data)
}

func TestOpenSearchClient_DeleteIndex(t *testing.T) {
	OsClient.DeleteIndex("image")
}
