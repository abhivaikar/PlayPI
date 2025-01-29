package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MessagePayload struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Message  string `json:"message,omitempty"`
	To       string `json:"to,omitempty"`
}

func sendRequest(t *testing.T, method, url string, payload interface{}) map[string]interface{} {
	var req *http.Request
	var err error

	// Prepare the request based on the method
	if method == http.MethodPost {
		var body []byte
		body, err = json.Marshal(payload)
		require.NoError(t, err)

		log.Printf("Sending payload: %s", string(body))

		req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
	} else if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
	} else {
		t.Fatalf("Unsupported HTTP method: %s", method)
	}

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read and parse the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP 200 OK, got %d", resp.StatusCode)

	var result map[string]interface{}
	err = json.Unmarshal(responseBody, &result)
	require.NoError(t, err)

	return result
}

func TestLiveChatE2E(t *testing.T) {
	clientAURL := "http://localhost:8090"
	clientBURL := "http://localhost:8091"

	// Step 1: Client A connects
	connectPayload := map[string]string{"websocket_server_url": "ws://localhost:8086/ws"}
	connectAResp := sendRequest(t, http.MethodPost, clientAURL+"/connect", connectPayload)
	require.NotNil(t, connectAResp)
	require.Contains(t, connectAResp, "message")
	require.Equal(t, "Connected successfully", connectAResp["message"])
	log.Println("Client A connected successfully")

	// Step 2: Client A reads the "you have joined" message
	time.Sleep(2 * time.Second)
	readAResp := sendRequest(t, http.MethodGet, clientAURL+"/read", nil)
	require.NotNil(t, readAResp)
	require.Contains(t, readAResp["message"], "You have connected as")
	log.Println("Client A read the welcome message after joining")

	// Extract username from the message
	message := readAResp["message"].(string)
	parts := strings.Split(message, "You have connected as ")
	require.Len(t, parts, 2, "Expected message format: 'You have connected as <username>'")
	clientAUsername := parts[1]

	// Step 3: Client B connects
	time.Sleep(2 * time.Second)
	connectBResp := sendRequest(t, http.MethodPost, clientBURL+"/connect", connectPayload)
	require.NotNil(t, connectBResp)
	require.Contains(t, connectBResp, "message")
	require.Equal(t, "Connected successfully", connectBResp["message"])
	log.Println("Client B connected successfully")

	// Step 4: Client A reads the message about Client B joining
	time.Sleep(2 * time.Second)
	readAResp = sendRequest(t, http.MethodGet, clientAURL+"/read", nil)
	require.NotNil(t, readAResp)
	require.Contains(t, readAResp["message"], "has joined the chat.")
	log.Println("Client A read the message about Client B joining")

	message = readAResp["message"].(string)
	parts = strings.Split(message, " has joined the chat")
	require.Len(t, parts, 2, "Expected message format: '<username> has joined the chat'")
	clientBUsername := parts[0]

	// Step 5: Client A sends a public message
	time.Sleep(2 * time.Second)
	sendAPublicMessageResp := sendRequest(t, http.MethodPost, clientAURL+"/send", MessagePayload{
		Type:     "chat",
		Username: clientAUsername,
		Message:  "Hello everyone!",
	})
	require.NotNil(t, sendAPublicMessageResp)
	require.Equal(t, "Message sent successfully", sendAPublicMessageResp["message"])
	log.Println("Client A sent a public message")

	// Step 6: Client B reads the public message
	time.Sleep(2 * time.Second)
	readBResp := sendRequest(t, http.MethodGet, clientBURL+"/read", nil)
	require.NotNil(t, readBResp)
	require.Contains(t, readBResp["message"], "Hello everyone!")
	log.Println("Client B read the public message")

	// Step 7: Client A sends a private message to Client B
	time.Sleep(2 * time.Second)
	sendAPrivateMessageResp := sendRequest(t, http.MethodPost, clientAURL+"/send", MessagePayload{
		Type:     "private",
		Username: clientAUsername,
		Message:  "Hello B!",
		To:       clientBUsername,
	})
	require.NotNil(t, sendAPrivateMessageResp)
	require.Equal(t, "Message sent successfully", sendAPrivateMessageResp["message"])
	log.Println("Client A sent a private message to Client B")

	// Step 8: Client B reads the private message
	time.Sleep(2 * time.Second)
	readBResp = sendRequest(t, http.MethodGet, clientBURL+"/read", nil)
	require.Contains(t, readBResp["message"], "Hello B!")
	log.Println("Client B read the private message from Client A")

	// Step 9: Client A disconnects
	time.Sleep(2 * time.Second)
	disconnectAResp := sendRequest(t, http.MethodPost, clientAURL+"/disconnect", nil)
	require.NotNil(t, disconnectAResp)
	require.Contains(t, disconnectAResp, "message")
	require.Equal(t, "Disconnected successfully", disconnectAResp["message"])
	log.Println("Client A disconnected successfully")

	// Step 10: Client B reads the disconnection message
	time.Sleep(2 * time.Second)
	readBResp = sendRequest(t, http.MethodGet, clientBURL+"/read", nil)
	require.Contains(t, readBResp["message"], clientAUsername+" has left the chat.")
	log.Println("Client B read the disconnection message from Client A")
}
