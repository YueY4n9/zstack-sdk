package common

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	accountName    string
	password       string
	serverEndpoint string
	sessionID      string
	httpClient     *http.Client
}

func (client *Client) Init(AccountName, Password, ServerEndpoint string) error {
	hsha512 := sha512.New()
	_, _ = io.WriteString(hsha512, Password)
	client.accountName = AccountName
	client.password = fmt.Sprintf("%x", hsha512.Sum(nil))
	client.serverEndpoint = ServerEndpoint
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.httpClient = &http.Client{Transport: tr}
	return client.login()
}

func (client *Client) login() error {
	login := LoginRequest{
		LoginContent: map[string]string{
			"password":    client.password,
			"accountName": client.accountName,
		},
		Tags: Tags{
			SystemTags: []string{},
			UserTags:   []string{},
		},
	}
	requestBody, _ := json.Marshal(login)
	resp, err := client.CreateRequestWithURI(http.MethodPut, loginURI, requestBody)
	if err != nil {
		return errors.Wrap(err, "Get error while login request")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Get error while getting data from login response")
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, respBody, "", "    ")
	log.Println(string(prettyJSON.Bytes()))

	loginResponse := LoginResponse{}
	errorResponse := ErrorResponse{}

	if int(resp.StatusCode/100) != 2 {
		if err := json.Unmarshal(respBody, &errorResponse); err != nil {
			return errors.Wrap(err, "Get error while decoding login response")
		}
		return errors.New(errorResponse.Error.Description + " " + errorResponse.Error.Details)
	} else {
		if err := json.Unmarshal(respBody, &loginResponse); err != nil {
			return errors.Wrap(err, "Get error while decoding login response")
		}
	}

	//if loginResponse.Error != nil {
	//	return errors.New(loginResponse.Error.Description + " " + loginResponse.Error.Details)
	//}

	client.sessionID = loginResponse.Inventory.UUID
	return nil
}

func (client *Client) captcha() {

}

func (client *Client) CreateRequestWithURI(method, uri string, body []byte) (*http.Response, error) {
	urlPath := client.serverEndpoint + uri
	httpRequest, err := http.NewRequest(method, urlPath, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Add("Authorization", "OAuth "+client.sessionID)
	resp, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
