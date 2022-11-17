package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/env"
	"log"
	"net/http"
)

var EsClient *ESClient

type ESClient struct {
	*elasticsearch.Client
}

func init() {
	//address := env.GetEnv(env.OpenSearchAddress, "https://10.40.101.1:31367")
	address := env.GetEnv(env.OpenSearchAddress, "https://10.40.101.1:9200")
	username := env.GetEnv(env.OpenSearchUsername, "admin")
	password := env.GetEnv(env.OpenSearchAddress, "admin")
	cfg := elasticsearch.Config{
		Addresses: []string{
			address,
		},
		Username: username,
		Password: password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("new es client err: %v", err))
	}

	//res, err := es.Info()
	//if err != nil {
	//	log.Fatalf("Error getting response: %s", err)
	//}
	//log.Println(res)

	EsClient = &ESClient{Client: es}
}

func (e *ESClient) IndexRequest(documentId string, body []byte) {
	req := esapi.IndexRequest{
		Index:      env.OpenSearchIndex,
		DocumentID: documentId,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), e.Client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 1)
	} else {
		log.Printf("res: %s", res.String())
	}

	return
}
