package incloud

import (
	"crypto/tls"
// 	"encoding/json"
// 	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
// 	"strconv"
// 	"strings"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
    tenant_id  string
    api_token  string
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}

// NewClient -
func NewClient(host, username, password *string, token *string) (*Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	c := Client{
	    HostURL: *host,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
			Transport: tr,
		},
		
	}

	if host != nil {
		c.HostURL = *host
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
    
	req.Header.Set("Authorization", c.api_token)
	

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}