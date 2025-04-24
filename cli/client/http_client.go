package client

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// sendRequest is a helper function to send an HTTP request
func sendRequest(method, addr, path, data string) (string, error) {
	var req *http.Request
	var err error

	if method == "POST" || method == "DELETE" {
		req, err = http.NewRequest(method, fmt.Sprintf("http://%s/%s", addr, path), strings.NewReader(data))
		if err != nil {
			return "", fmt.Errorf("error creating %s request: %v", method, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("http://%s/%s?%s", addr, path, data), nil)
		if err != nil {
			return "", fmt.Errorf("error creating %s request: %v", method, err)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error performing %s request: %v", method, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	return string(body), nil
}

// Get sends a GET request to the specified address with the given key.
func Get(addr, key string) (string, error) {
	return sendRequest("GET", addr, "get", "key="+key)
}

// Set sends a POST request to the specified address with the given key-value pair.
func Set(addr, key, value string) (string, error) {
	data := fmt.Sprintf("key=%s&value=%s", key, value)

	return sendRequest("POST", addr, "set", data)
}

// Delete sends a DELETE request to the specified address to delete the value associated with the given key.
func Delete(addr, key string) (string, error) {
	data := fmt.Sprintf("key=%s", key)

	return sendRequest("DELETE", addr, "delete", data)
}
