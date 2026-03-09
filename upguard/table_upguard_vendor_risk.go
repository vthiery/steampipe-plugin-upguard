package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// VendorRisk represents a risk detected for a vendor.
type VendorRisk struct {
	RiskID        string                 `json:"risk_id"`
	Severity      string                 `json:"severity"`
	Category      string                 `json:"category"`
	DetectedAt    string                 `json:"detected_at"`
	LastScannedAt string                 `json:"last_scanned_at"`
	Hostnames     []string               `json:"hostnames"`
	Sources       interface{}            `json:"sources"`
	Meta          map[string]interface{} `json:"meta"`
}

// VendorRiskSource represents the source of a vendor risk.
type VendorRiskSource struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

// vendorRisksResponse is the envelope returned by GET /risks/vendors.
type vendorRisksResponse struct {
	Risks []VendorRisk `json:"risks"`
}

//// TABLE DEFINITION

func tableUpguardVendorRisk() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_vendor_risk",
		Description: "List active risks detected for a specific vendor in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listVendorRisks,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vendor_primary_hostname", Require: plugin.Required},
			},
		},
		Columns: []*plugin.Column{
			{Name: "vendor_primary_hostname", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "Primary hostname of the vendor (must be provided in WHERE clause)."},
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

func listVendorRisks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// Vendor primary hostname is required
	vendorHostname := d.EqualsQualString("vendor_primary_hostname")
	if vendorHostname == "" {
		return nil, fmt.Errorf("vendor_primary_hostname must be specified in WHERE clause")
	}

	params := map[string]string{
		"primary_hostname": vendorHostname,
		"include_sources":  "true",
		"include_meta":     "true",
	}

	// Add min_severity filter if specified
	if d.EqualsQuals["min_severity"] != nil {
		params["min_severity"] = d.EqualsQualString("min_severity")
	}

	var result vendorRisksResponse
	if err := client.get(ctx, "/risks/vendors", params, &result); err != nil {
		return nil, fmt.Errorf("listing vendor risks: %w", err)
	}

	for _, risk := range result.Risks {
		d.StreamListItem(ctx, risk)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
