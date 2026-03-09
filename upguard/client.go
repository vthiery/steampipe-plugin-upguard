package upguard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const baseURL = "https://cyber-risk.upguard.com/api/public"

// Client wraps the UpGuard REST API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// newClient creates a new Client using the given API key.
func newClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getClient retrieves a configured Client from the plugin connection config.
func getClient(ctx context.Context, d *plugin.QueryData) (*Client, error) {
	cfg := GetConfig(d.Connection)
	if cfg.APIKey == nil {
		return nil, fmt.Errorf("api_key must be configured for the upguard plugin")
	}
	return newClient(*cfg.APIKey), nil
}

// get performs an authenticated GET request and unmarshals the response body into result.
func (c *Client) get(ctx context.Context, path string, params map[string]string, result interface{}) error {
	reqURL := baseURL + path
	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		reqURL += "?" + values.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// Shared sub-types used across multiple tables
// ---------------------------------------------------------------------------

// PaginationMeta holds pagination info returned by endpoints.
type PaginationMeta struct {
	NextPageToken string `json:"next_page_token"`
	TotalResults  int    `json:"total_results"`
	PageSize      int    `json:"page_size"`
}

// CategoryScores represents the breakdown of scores by category.
type CategoryScores struct {
	AttackSurface           int `json:"attackSurface"`
	BrandProtection         int `json:"brandProtection"`
	BrandReputation         int `json:"brandReputation"`
	DataLeakage             int `json:"dataLeakage"`
	DNS                     int `json:"dns"`
	EmailSecurity           int `json:"emailSecurity"`
	Encryption              int `json:"encryption"`
	IPDomainReputation      int `json:"ipDomainReputation"`
	NetworkSecurity         int `json:"networkSecurity"`
	OperationalRisk         int `json:"operationalRisk"`
	Phishing                int `json:"phishing"`
	Questionnaires          int `json:"questionnaires"`
	VulnerabilityManagement int `json:"vulnerabilityManagement"`
	WebsiteSecurity         int `json:"websiteSecurity"`
}

// RiskCounts represents the count of risks by severity.
type RiskCounts struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}
