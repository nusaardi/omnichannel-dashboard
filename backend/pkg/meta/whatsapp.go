package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const whatsappAPIURL = "https://graph.facebook.com/v19.0"

// WhatsAppClient handles WhatsApp Cloud API interactions
type WhatsAppClient struct {
	accessToken string
	phoneID     string
	businessID  string
	httpClient  *http.Client
}

// NewWhatsAppClient creates a new WhatsApp API client
func NewWhatsAppClient(accessToken, phoneID, businessID string) *WhatsAppClient {
	return &WhatsAppClient{
		accessToken: accessToken,
		phoneID:     phoneID,
		businessID:  businessID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TextMessage represents a text message to send
type TextMessage struct {
	To   string `json:"to"`
	Type string `json:"type"`
	Text struct {
		Body string `json:"body"`
	} `json:"text"`
	MessagingProduct string `json:"messaging_product"`
}

// SendTextResponse represents the API response
type SendTextResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

// SendText sends a text message via WhatsApp
func (c *WhatsAppClient) SendText(to, message string) (*SendTextResponse, error) {
	payload := TextMessage{
		To:               to,
		Type:             "text",
		MessagingProduct: "whatsapp",
	}
	payload.Text.Body = message

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", whatsappAPIURL, c.phoneID)
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

	var result SendTextResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// TemplateMessage represents a template message
type TemplateMessage struct {
	To               string `json:"to"`
	Type             string `json:"type"`
	MessagingProduct string `json:"messaging_product"`
	Template         struct {
		Name     string `json:"name"`
		Language struct {
			Code string `json:"code"`
		} `json:"language"`
		Components []TemplateComponent `json:"components,omitempty"`
	} `json:"template"`
}

type TemplateComponent struct {
	Type       string              `json:"type"`
	Parameters []TemplateParameter `json:"parameters,omitempty"`
}

type TemplateParameter struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// SendTemplate sends a template message
func (c *WhatsAppClient) SendTemplate(to, templateName, languageCode string, params []string) (*SendTextResponse, error) {
	msg := TemplateMessage{
		To:               to,
		Type:             "template",
		MessagingProduct: "whatsapp",
	}
	msg.Template.Name = templateName
	msg.Template.Language.Code = languageCode

	if len(params) > 0 {
		var parameters []TemplateParameter
		for _, p := range params {
			parameters = append(parameters, TemplateParameter{
				Type: "text",
				Text: p,
			})
		}
		msg.Template.Components = []TemplateComponent{
			{
				Type:       "body",
				Parameters: parameters,
			},
		}
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal template: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", whatsappAPIURL, c.phoneID)
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
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var result SendTextResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// MarkAsRead marks a message as read
func (c *WhatsAppClient) MarkAsRead(messageID string) error {
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s/messages", whatsappAPIURL, c.phoneID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}
