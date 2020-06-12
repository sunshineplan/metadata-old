package metadata

import (
	"encoding/json"
	"errors"
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

var debug = false

// SetDebug sets the debug status
// Setting this to true causes the panics to be thrown and logged onto the console.
// Setting this to false causes the errors to be saved in the Error field in the returned struct.
func SetDebug(d bool) {
	debug = d
}

// Get metadata from metadata server
func Get(metadata string, c Config) (interface{}, error) {
	var value interface{}
	var result []byte
	client := &http.Client{}
	server := c.Server + "/" + metadata
	req, err := http.NewRequest("GET", server, nil)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + server)
		}
		return "", errors.New("error creating get request to " + server)
	}
	req.Header.Add(c.VerifyHeader, c.VerifyValue)
	resp, err := client.Do(req)
	if err != nil {
		if debug {
			panic("Couldn't perform GET request to " + server)
		}
		return "", errors.New("couldn't perform GET request to " + server)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if debug {
				panic("Couldn't read response from " + server)
			}
			return "", errors.New("Couldn't read response from " + server)
		}
		result = bodyBytes
	} else {
		if debug {
			panic("No StatusOK response from " + server)
		}
		return "", errors.New("No StatusOK response from " + server)
	}
	err = json.Unmarshal(result, &value)
	if err != nil {
		if debug {
			panic("JSON Decode Error")
		}
		return "", err
	}
	return value, nil
}
