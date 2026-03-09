package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// VendorIP represents an IP address belonging to a vendor.
type VendorIP struct {
	IP             string   `json:"ip"`
	Owner          string   `json:"owner"`
	Country        string   `json:"country"`
	ASN            int      `json:"asn"`
	ASName         string   `json:"as_name"`
	AutomatedScore int      `json:"automated_score"`
	Labels         []string `json:"labels"`
}

// vendorIPsResponse is the envelope returned by GET /vendor/ips.
type vendorIPsResponse struct {
	IPs           []VendorIP `json:"ips"`
	NextPageToken string     `json:"next_page_token"`
	TotalResults  int        `json:"total_results"`
}

//// TABLE DEFINITION

func tableUpguardVendorIP() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_vendor_ip",
		Description: "List IP addresses for a specific vendor in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listVendorIPs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "vendor_primary_hostname", Require: plugin.Required},
				{Name: "labels", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"vendor_primary_hostname", "ip"}),
			Hydrate:    getVendorIP,
		},
		Columns: []*plugin.Column{
			{Name: "vendor_primary_hostname", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "Primary hostname of the vendor (required)."},
			{Name: "ip", Type: proto.ColumnType_IPADDR, Description: "The IP address."},
			{Name: "owner", Type: proto.ColumnType_STRING, Description: "Owner of the IP address."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country where the IP is located."},
			{Name: "asn", Type: proto.ColumnType_INT, Description: "Autonomous System Number."},
			{Name: "as_name", Type: proto.ColumnType_STRING, Description: "Autonomous System Name."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score for the IP."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels assigned to the IP."},
			{Name: "services", Type: proto.ColumnType_JSON, Description: "Services detected on the IP."},
			{Name: "check_results", Type: proto.ColumnType_JSON, Description: "Security check results for the IP."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listVendorIPs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

		var result vendorIPsResponse
		if err := client.get(ctx, "/vendor/ips", params, &result); err != nil {
			return nil, fmt.Errorf("listing vendor IPs: %w", err)
		}

		for _, ip := range result.IPs {
			d.StreamListItem(ctx, ip)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextPageToken == "" || len(result.IPs) == 0 {
			break
		}
		pageToken = result.NextPageToken
	}

	return nil, nil
}

func getVendorIP(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	vendorHostname := d.EqualsQualString("vendor_primary_hostname")
	ip := d.EqualsQualString("ip")

	if vendorHostname == "" || ip == "" {
		return nil, nil
	}

	params := map[string]string{
		"vendor_primary_hostname": vendorHostname,
		"ip":                      ip,
	}

	var result IPDetails
	if err := client.get(ctx, "/vendor/ip", params, &result); err != nil {
		return nil, fmt.Errorf("getting vendor IP %s: %w", ip, err)
	}

	return result, nil
}
