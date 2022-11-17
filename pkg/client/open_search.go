package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/env"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/logging"
	"io"
	"net/http"
	"net/url"
)

var OsClient *OpenSearchClient

type OpenSearchClient struct {
	http      *http.Client
	osAddress string
	username  string
	password  string
	index     string
	baseURL   *url.URL
}

func init() {
	address := env.GetEnv(env.OpenSearchAddress)
	username := env.GetEnv(env.OpenSearchUsername)
	password := env.GetEnv(env.OpenSearchPassword)
	index := env.GetEnv(env.OpenSearchIndex)

	if address == "" {
		panic(fmt.Sprintf("please set env: %s", env.OpenSearchAddress))
	}

	if username == "" {
		panic(fmt.Sprintf("please set env: %s", env.OpenSearchUsername))
	}

	if password == "" {
		panic(fmt.Sprintf("please set env: %s", env.OpenSearchPassword))
	}

	if index == "" {
		panic(fmt.Sprintf("please set env: %s", env.OpenSearchIndex))
	}

	baseUrl, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprintf("openSearch address is err: %v", err))
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	OsClient = &OpenSearchClient{
		http:      httpClient,
		osAddress: address,
		username:  username,
		password:  password,
		index:     index,
		baseURL:   baseUrl,
	}
}

func (o *OpenSearchClient) Info() (info string) {
	var bodyBytes []byte
	var err error

	request, err := http.NewRequest(
		"GET",
		o.osAddress,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		logging.GetLogger().Error().Msgf("new request err: %v", err)
		return
	}
	request.SetBasicAuth(o.username, o.password)

	response, err := o.http.Do(request)
	if err != nil {
		logging.GetLogger().Error().Msgf("send http request err: %v", err)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.GetLogger().Error().Msgf("get openSearch info err: %v", err)
		return
	}

	logging.GetLogger().Info().Msgf("openSearch info: %s", string(body))
	info = string(body)

	return
}

func (o *OpenSearchClient) Add(documentId string, requestBody []byte) (err error) {
	var (
		body []byte
	)

	path := fmt.Sprintf("%s/_doc/%s", o.index, documentId)
	request, err := http.NewRequest(
		"PUT",
		o.url(path).String(),
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		logging.GetLogger().Error().Msgf("new request err: %v", err)
		return
	}
	request.SetBasicAuth(o.username, o.password)
	request.Header.Set("Content-Type", "application/json")

	response, err := o.http.Do(request)
	if err != nil {
		logging.GetLogger().Error().Msgf("send add http request err: %v", err)
		return
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		logging.GetLogger().Error().Msgf("openSearch add err: %v", err)
		return
	}

	logging.GetLogger().Info().Msgf("openSearch add response: %s", string(body))
	return
}

func (o *OpenSearchClient) DeleteIndex(index string) {
	var bodyBytes []byte
	var err error

	request, err := http.NewRequest(
		"DELETE",
		o.url(index).String(),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		logging.GetLogger().Error().Msgf("new request err: %v", err)
		return
	}
	request.SetBasicAuth(o.username, o.password)

	response, err := o.http.Do(request)
	if err != nil {
		logging.GetLogger().Error().Msgf("send http request err: %v", err)
		return
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logging.GetLogger().Error().Msgf("delete openSearch index err: %v", err)
		return
	}

	logging.GetLogger().Info().Msgf("delete openSearch index response: %s", string(body))

	return
}

func (o *OpenSearchClient) url(path string) *url.URL {
	u := *o.baseURL
	u.Path = path
	return &u
}
