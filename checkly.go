package checkly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// NewClient takes a Checkly API key, and returns a Client ready to use.
func NewClient(apiKey string) Client {
	return Client{
		apiKey:     apiKey,
		URL:        getEnv("CHECKLY_API_URL", "https://api.checklyhq.com"),
		HTTPClient: http.DefaultClient,
	}
}

// Create creates a new check with the specified details. It returns the
// check ID of the newly-created check, or an error.
func (c *Client) Create(check Check) (string, error) {
	data, err := json.Marshal(check)
	if err != nil {
		return "", err
	}
	status, res, err := c.MakeAPICall(http.MethodPost, "checks", data)
	if err != nil {
		return "", err
	}
	if status != http.StatusCreated {
		return "", fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	var result Check
	if err = json.NewDecoder(strings.NewReader(res)).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return result.ID, nil
}

// Update updates an existing check with the specified details. It returns a
// non-nil error if the request failed.
func (c *Client) Update(ID string, check Check) error {
	data, err := json.Marshal(check)
	if err != nil {
		return err
	}
	status, res, err := c.MakeAPICall(http.MethodPut, "checks/"+ID, data)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	var result Check
	if err = json.NewDecoder(strings.NewReader(res)).Decode(&result); err != nil {
		return fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return nil
}

// Delete deletes the check with the specified ID. It returns a non-nil
// error if the request failed.
func (c *Client) Delete(ID string) error {
	status, res, err := c.MakeAPICall(http.MethodDelete, "checks/"+ID, nil)
	if err != nil {
		return err
	}
	if status != http.StatusNoContent {
		return fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	return nil
}

// Get takes the ID of an existing check, and returns the check parameters, or
// an error.
func (c *Client) Get(ID string) (Check, error) {
	status, res, err := c.MakeAPICall(http.MethodGet, "checks/"+ID, nil)
	if err != nil {
		return Check{}, err
	}
	if status != http.StatusOK {
		return Check{}, fmt.Errorf("unexpected response status %d: %q", status, res)
	}
	check := Check{}
	if err = json.NewDecoder(strings.NewReader(res)).Decode(&check); err != nil {
		return Check{}, fmt.Errorf("decoding error for data %s: %v", res, err)
	}
	return check, nil
}

// MakeAPICall calls the Checkly API with the specified URL and data, and
// returns the HTTP status code and string data of the response.
func (c *Client) MakeAPICall(method string, URL string, data []byte) (statusCode int, response string, err error) {
	requestURL := c.URL + "/v1/" + URL
	req, err := http.NewRequest(method, requestURL, bytes.NewBuffer(data))
	if err != nil {
		return 0, "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("content-type", "application/json")
	if c.Debug != nil {
		requestDump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			return 0, "", fmt.Errorf("error dumping HTTP request: %v", err)
		}
		fmt.Fprintln(c.Debug, string(requestDump))
		fmt.Fprintln(c.Debug)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	if c.Debug != nil {
		c.dumpResponse(resp)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", err
	}
	return resp.StatusCode, string(res), nil
}

// dumpResponse writes the raw response data to the debug output, if set, or
// standard error otherwise.
func (c *Client) dumpResponse(resp *http.Response) {
	// ignore errors dumping response - no recovery from this
	responseDump, _ := httputil.DumpResponse(resp, true)
	fmt.Fprintln(c.Debug, string(responseDump))
	fmt.Fprintln(c.Debug)
}
