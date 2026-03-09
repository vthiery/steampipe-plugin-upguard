package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// OrganisationRisk represents a risk detected for the organization.
type OrganisationRisk struct {
	RiskID        string                 `json:"risk_id"`
	Severity      string                 `json:"severity"`
	Category      string                 `json:"category"`
	DetectedAt    string                 `json:"detected_at"`
	LastScannedAt string                 `json:"last_scanned_at"`
	Hostnames     []string               `json:"hostnames"`
	Sources       interface{}            `json:"sources"`
	Meta          map[string]interface{} `json:"meta"`
}

// organisationRisksResponse is the envelope returned by GET /risks.
type organisationRisksResponse struct {
	Risks []OrganisationRisk `json:"risks"`
}

//// TABLE DEFINITION

func tableUpguardOrganisationRisk() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_organisation_risk",
		Description: "List active risks detected for your organization in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listOrganisationRisks,
		},
		Columns: []*plugin.Column{
			{Name: "risk_id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the risk type."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Severity level (info, low, medium, high, critical)."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Risk category."},
			{Name: "detected_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the risk was first detected."},
			{Name: "last_scanned_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the risk was last scanned."},
			{Name: "hostnames", Type: proto.ColumnType_JSON, Description: "Hostnames where the risk was detected."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "Source details including hostname, IP, and port."},
			{Name: "meta", Type: proto.ColumnType_JSON, Description: "Additional metadata about the risk."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listOrganisationRisks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"include_sources": "true",
		"include_meta":    "true",
		"min_severity":    "info",
	}

	// Add min_severity filter if specified
	if d.EqualsQuals["min_severity"] != nil {
		params["min_severity"] = d.EqualsQualString("min_severity")
	}

	var result organisationRisksResponse
	if err := client.get(ctx, "/risks", params, &result); err != nil {
		return nil, fmt.Errorf("listing organisation risks: %w", err)
	}

	for _, risk := range result.Risks {
		d.StreamListItem(ctx, risk)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
