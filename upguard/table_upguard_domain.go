package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// DomainListItem represents a domain in the LIST /domains response.
// The LIST endpoint only returns minimal fields.
type DomainListItem struct {
	Hostname      string `json:"hostname"`
	Active        bool   `json:"active"`
	PrimaryDomain bool   `json:"primary_domain"`
}

// domainsResponse is the envelope returned by GET /domains.
type domainsResponse struct {
	Domains       []DomainListItem `json:"domains"`
	NextPageToken string           `json:"next_page_token"`
	TotalResults  int              `json:"total_results"`
}

// Domain represents detailed information about a domain from GET /domain.
// This struct contains all fields available from the detailed domain endpoint.
type Domain struct {
	Hostname           string        `json:"hostname"`
	Active             bool          `json:"active"`
	AutomatedScore     int           `json:"automated_score"`
	ScannedAt          string        `json:"scanned_at"`
	ARecords           []string      `json:"a_records"`
	Labels             []string      `json:"labels"`
	CheckResults       []CheckResult `json:"check_results"`
	WaivedCheckResults []CheckResult `json:"waived_check_results"`
}

// CheckResult represents a security check result.
// Fields match the actual API response structure.
type CheckResult struct {
	ID           string                   `json:"id"`
	RiskType     string                   `json:"riskType"`
	Pass         bool                     `json:"pass"`
	Severity     int                      `json:"severity"`      // Numeric severity level
	SeverityName string                   `json:"severityName"`  // String severity ("info", "low", etc.)
	Category     string                   `json:"category"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description"`
	CheckedAt    string                   `json:"checked_at"`
	Sources      []string                 `json:"sources"`
	Expected     []map[string]interface{} `json:"expected"`
	Actual       []map[string]interface{} `json:"actual"`
}

//// TABLE DEFINITION

func tableUpguardDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_domain",
		Description: "List and inspect domains in your UpGuard account.",
		List: &plugin.ListConfig{
			Hydrate: listDomains,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "active", Require: plugin.Optional},
				{Name: "labels", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("hostname"),
			Hydrate:    getDomain,
		},
		Columns: []*plugin.Column{
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "The domain hostname."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the domain is active."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score for the domain."},
			{Name: "scanned_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ScannedAt").Transform(transform.NullIfZeroValue), Description: "When the domain was last scanned."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels assigned to the domain."},
			{Name: "a_records", Type: proto.ColumnType_JSON, Description: "DNS A records for the domain."},
			{Name: "check_results", Type: proto.ColumnType_JSON, Description: "Security check results for the domain."},
			{Name: "waived_check_results", Type: proto.ColumnType_JSON, Description: "Waived security check results."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	params := map[string]string{
		"page_size": "1000",
		"active":    "true",
		"inactive":  "true",
	}

	// Add filters if specified
	if d.EqualsQuals["active"] != nil {
		if d.EqualsQuals["active"].GetBoolValue() {
			params["active"] = "true"
			params["inactive"] = "false"
		} else {
			params["active"] = "false"
			params["inactive"] = "true"
		}
	}

	if d.EqualsQuals["labels"] != nil {
		labels := d.EqualsQuals["labels"].GetJsonbValue()
		if labels != "" {
			params["labels"] = labels
		}
	}

	for {
		if pageToken != "" {
			params["page_token"] = pageToken
		}

		var result domainsResponse
		if err := client.get(ctx, "/domains", params, &result); err != nil {
			return nil, fmt.Errorf("listing domains: %w", err)
		}

		for _, listItem := range result.Domains {
			// Convert DomainListItem to Domain for consistent column access
			domain := Domain{
				Hostname: listItem.Hostname,
				Active:   listItem.Active,
				// AutomatedScore, ScannedAt, etc. not available from LIST endpoint
			}
			d.StreamListItem(ctx, domain)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextPageToken == "" || len(result.Domains) == 0 {
			break
		}
		pageToken = result.NextPageToken
	}

	return nil, nil
}

func getDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	hostname := d.EqualsQualString("hostname")
	if hostname == "" {
		return nil, nil
	}

	params := map[string]string{
		"hostname": hostname,
	}

	var result Domain
	if err := client.get(ctx, "/domain", params, &result); err != nil {
		return nil, fmt.Errorf("getting domain %s: %w", hostname, err)
	}

	return result, nil
}
