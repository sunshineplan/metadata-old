package metadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config contain metadata server address and verify header
type Config struct {
	// metadata server address
	Server string
	// metadata server verify header name
	VerifyHeader string
	// metadata server verify header value
	VerifyValue string
}

// Get metadata from metadata server
func (c *Config) Get(metadata string) ([]byte, error) {
	return c.GetWithClient(metadata, &http.Client{Transport: &http.Transport{Proxy: nil}})
}

// GetWithClient metadata from metadata server with custom http.Client
func (c *Config) GetWithClient(metadata string, client *http.Client) ([]byte, error) {
	server := c.Server + "/" + metadata
	req, err := http.NewRequest("GET", server, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to %s: %v", server, err)
	}
	req.Header.Add(c.VerifyHeader, c.VerifyValue)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request to %s: %v", server, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("No StatusOK response from %s", server)
	}
	return ioutil.ReadAll(resp.Body)
}
