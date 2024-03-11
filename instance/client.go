package instance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YueY4n9/zstack-sdk/common"
	"github.com/sirupsen/logrus"
)

type Client struct {
	common.Client
}

func NewInstanceClient(accountName string, password string, endpoint string) *Client {
	client := &Client{}
	err := client.Init(accountName, password, endpoint)
	if err != nil {
		logrus.Errorf("Get instance client error, reason: %s", err.Error())
		return nil
	}
	return client
}

func (c *Client) QueryInstances() ([]*VMInstanceInventory, error) {
	resp, err := c.CreateRequestWithURI(http.MethodGet, queryInstancesURI, nil)
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	responseStruct := QueryInstanceResponse{}
	if err = json.Unmarshal(responseBody, &responseStruct); err != nil {
		logrus.Warnf("Unmarshaling response when Querying instance. Error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		if responseStruct.Error != nil {
			return nil, responseStruct.Error.WrapError()
		}
		return nil, fmt.Errorf("status code %d,Error massage %s", resp.StatusCode, string(responseBody))
	}
	return responseStruct.Inventories, nil
}
