package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const instagramAPIURL = "https://graph.facebook.com/v19.0"

// InstagramClient handles Instagram Graph API interactions for messaging
type InstagramClient struct {
	accessToken string
	accountID   string
	httpClient  *http.Client
}

// NewInstagramClient creates a new Instagram API client
func NewInstagramClient(accessToken, accountID string) *InstagramClient {
	return &InstagramClient{
		accessToken: accessToken,
		accountID:   accountID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IGMessage represents an Instagram DM message
type IGMessage struct {
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text,omitempty"`
	} `json:"message"`
}

// IGMessageResponse represents the API response
type IGMessageResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_id"`
}

// SendText sends a text message via Instagram DM
func (c *InstagramClient) SendText(recipientID, message string) (*IGMessageResponse, error) {
	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": message,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", instagramAPIURL, c.accountID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var result IGMessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetConversations retrieves Instagram DM conversations
func (c *InstagramClient) GetConversations(limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s/conversations?fields=participants,messages{message,from,created_time}&limit=%d",
		instagramAPIURL, c.accountID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// GetUserProfile retrieves Instagram user profile
func (c *InstagramClient) GetUserProfile(userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s?fields=id,username,name,profile_picture_url", instagramAPIURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}
