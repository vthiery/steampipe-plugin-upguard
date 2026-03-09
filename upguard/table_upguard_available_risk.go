package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// AvailableRisk represents a risk type available in the UpGuard platform.
type AvailableRisk struct {
	ID          string `json:"id"`
	Risk        string `json:"risk"`
	Finding     string `json:"finding"`
	RiskDetails string `json:"riskDetails"`
	Remediation string `json:"remediation"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Group       string `json:"group"`
	Severity    string `json:"severity"`
	RiskType    string `json:"riskType"`
	RiskSubtype string `json:"riskSubtype"`
	Generic     bool   `json:"generic"`
}

// availableRisksResponse is the envelope returned by GET /available_risks/v2.
type availableRisksResponse []AvailableRisk

//// TABLE DEFINITION

func tableUpguardAvailableRisk() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_available_risk",
		Description: "List all available risk types in the UpGuard platform.",
		List: &plugin.ListConfig{
			Hydrate: listAvailableRisks,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getAvailableRisk,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Unique identifier for the risk type."},
			{Name: "risk", Type: proto.ColumnType_STRING, Description: "Short description of the risk."},
			{Name: "finding", Type: proto.ColumnType_STRING, Description: "What was found."},
			{Name: "risk_details", Type: proto.ColumnType_STRING, Description: "Detailed explanation of the risk."},
			{Name: "remediation", Type: proto.ColumnType_STRING, Description: "How to remediate the risk."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Full description of the risk."},
			{Name: "category", Type: proto.ColumnType_STRING, Description: "Risk category."},
			{Name: "group", Type: proto.ColumnType_STRING, Description: "Risk group."},
			{Name: "severity", Type: proto.ColumnType_STRING, Description: "Severity level (pass, info, low, medium, high, critical)."},
			{Name: "risk_type", Type: proto.ColumnType_STRING, Description: "Type of risk."},
			{Name: "risk_subtype", Type: proto.ColumnType_STRING, Description: "Subtype of risk."},
			{Name: "generic", Type: proto.ColumnType_BOOL, Description: "Whether this is a generic risk template."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listAvailableRisks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result availableRisksResponse
	if err := client.get(ctx, "/available_risks/v2", nil, &result); err != nil {
		return nil, fmt.Errorf("listing available risks: %w", err)
	}

	for _, risk := range result {
		d.StreamListItem(ctx, risk)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getAvailableRisk(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	riskID := d.EqualsQualString("id")
	if riskID == "" {
		return nil, nil
	}

	params := map[string]string{
		"risk_id": riskID,
	}

	var result AvailableRisk
	if err := client.get(ctx, "/available_risks/risk", params, &result); err != nil {
		return nil, fmt.Errorf("getting risk %s: %w", riskID, err)
	}

	return result, nil
}
