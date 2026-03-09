package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// VendorDomain represents a domain belonging to a vendor.
type VendorDomain struct {
	Hostname       string   `json:"hostname"`
	Active         bool     `json:"active"`
	AutomatedScore int      `json:"automated_score"`
	ScannedAt      string   `json:"scanned_at"`
	Labels         []string `json:"labels"`
}

// vendorDomainsResponse is the envelope returned by GET /vendor/domains.
type vendorDomainsResponse struct {
	Domains       []VendorDomain `json:"domains"`
	NextPageToken string         `json:"next_page_token"`
	TotalResults  int            `json:"total_results"`
}

//// TABLE DEFINITION

func tableUpguardVendorDomain() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_vendor_domain",
		Description: "List domains for a specific vendor in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listVendorDomains,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vendor_primary_hostname", Require: plugin.Required},
				{Name: "active", Require: plugin.Optional},
				{Name: "labels", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"vendor_primary_hostname", "hostname"}),
			Hydrate:    getVendorDomain,
		},
		Columns: []*plugin.Column{
			{Name: "vendor_primary_hostname", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "Primary hostname of the vendor (required)."},
			{Name: "hostname", Type: proto.ColumnType_STRING, Description: "The domain hostname."},
			{Name: "active", Type: proto.ColumnType_BOOL, Description: "Whether the domain is active."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score for the domain."},
			{Name: "scanned_at", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ScannedAt").Transform(transform.NullIfZeroValue), Description: "When the domain was last scanned."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels assigned to the domain."},
			{Name: "a_records", Type: proto.ColumnType_JSON, Description: "DNS A records for the domain."},
			{Name: "check_results", Type: proto.ColumnType_JSON, Description: "Security check results for the domain."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listVendorDomains(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	vendorHostname := d.EqualsQualString("vendor_primary_hostname")
	if vendorHostname == "" {
		return nil, fmt.Errorf("vendor_primary_hostname must be specified in WHERE clause")
	}

	pageToken := ""
	params := map[string]string{
		"vendor_primary_hostname": vendorHostname,
		"page_size":               "1000",
		"active":                  "true",
		"inactive":                "true",
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

		var result vendorDomainsResponse
		if err := client.get(ctx, "/vendor/domains", params, &result); err != nil {
			return nil, fmt.Errorf("listing vendor domains: %w", err)
		}

		for _, domain := range result.Domains {
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

func getVendorDomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	vendorHostname := d.EqualsQualString("vendor_primary_hostname")
	hostname := d.EqualsQualString("hostname")

	if vendorHostname == "" || hostname == "" {
		return nil, nil
	}

	params := map[string]string{
		"vendor_primary_hostname": vendorHostname,
		"hostname":                hostname,
	}

	var result Domain
	if err := client.get(ctx, "/vendor/domain", params, &result); err != nil {
		return nil, fmt.Errorf("getting vendor domain %s: %w", hostname, err)
	}

	return result, nil
}
