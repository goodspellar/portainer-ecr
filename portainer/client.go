package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	client   *http.Client
	URL      *url.URL
	username string
	password string
	apiToken string
}

// NewClient returns a Portainer API client
func NewClient(host string, user string, pass string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	URL, err := url.Parse(fmt.Sprintf("http://%s:9000", host))
	if err != nil {
		return nil, err
	}

	// Get API Token
	authURL := fmt.Sprintf("%s%s", URL, "/api/auth")
	values := map[string]string{"Username": user, "Password": pass}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]string
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	token := "Bearer " + data["jwt"]

	c := Client{
		client:   httpClient,
		URL:      URL,
		username: user,
		password: pass,
		apiToken: token,
	}

	return &c, nil
}

// CallAPI provides the execution method for the Portainer API
func (c *Client) CallAPI(method string, path string, body interface{}) (*http.Response, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	apiURL := fmt.Sprintf("%s%s", c.URL, path)

	req, err := http.NewRequest(method, apiURL, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.apiToken)

	return c.client.Do(req)
}
