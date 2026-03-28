package tailscale

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TailscaleClient interface {
	CreateAuthKey(opts CreateKeyOptions) (string, error)
}

type Client struct {
	APIKey  string
	Tailnet string
	BaseURL string
}

func NewClient(apiKey, tailnet string) *Client {
	return &Client{
		APIKey:  apiKey,
		Tailnet: tailnet,
		BaseURL: "https://api.tailscale.com/api/v2",
	}
}

func (c *Client) CreateAuthKey(opts CreateKeyOptions) (string, error) {
	url := fmt.Sprintf("%s/tailnet/%s/keys", c.BaseURL, c.Tailnet)

	reqBody := CreateKeyRequest{
		Capabilities: Capabilities{
			Devices: Devices{
				Create: DeviceCreate{
					Reusable:      opts.Reusable,
					Ephemeral:     opts.Ephemeral,
					Preauthorized: true,
					Tags:          opts.Tags,
				},
			},
		},
		Expiry: opts.Expiry,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(c.APIKey, "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tailscale API error: %d", resp.StatusCode)
	}

	var result CreateKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Key, nil
}
