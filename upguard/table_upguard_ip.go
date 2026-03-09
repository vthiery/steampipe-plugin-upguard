package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// IP represents an IP address in UpGuard.
type IP struct {
	IP             string   `json:"ip"`
	Owner          string   `json:"owner"`
	Country        string   `json:"country"`
	ASN            int      `json:"asn"`
	ASName         string   `json:"as_name"`
	AutomatedScore int      `json:"automated_score"`
	Labels         []string `json:"labels"`
}

// ipsResponse is the envelope returned by GET /ips.
type ipsResponse struct {
	IPs           []IP   `json:"ips"`
	NextPageToken string `json:"next_page_token"`
	TotalResults  int    `json:"total_results"`
}

// IPDetails represents detailed information about an IP address.
type IPDetails struct {
	IP                 string        `json:"ip"`
	Owner              string        `json:"owner"`
	Country            string        `json:"country"`
	ASN                int           `json:"asn"`
	ASName             string        `json:"as_name"`
	AutomatedScore     int           `json:"automated_score"`
	Labels             []string      `json:"labels"`
	Services           []string      `json:"services"`
	Sources            []string      `json:"sources"`
	CheckResults       []CheckResult `json:"check_results"`
	WaivedCheckResults []CheckResult `json:"waived_check_results"`
}

//// TABLE DEFINITION

func tableUpguardIP() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_ip",
		Description: "List and inspect IP addresses in your UpGuard account.",
		List: &plugin.ListConfig{
			Hydrate: listIPs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "labels", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("ip"),
			Hydrate:    getIP,
		},
		Columns: []*plugin.Column{
			{Name: "ip", Type: proto.ColumnType_IPADDR, Description: "The IP address."},
			{Name: "owner", Type: proto.ColumnType_STRING, Description: "Owner of the IP address."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country where the IP is located."},
			{Name: "asn", Type: proto.ColumnType_INT, Description: "Autonomous System Number."},
			{Name: "as_name", Type: proto.ColumnType_STRING, Description: "Autonomous System Name."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score for the IP."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels assigned to the IP."},
			{Name: "services", Type: proto.ColumnType_JSON, Description: "Services detected on the IP."},
			{Name: "sources", Type: proto.ColumnType_JSON, Description: "Sources where the IP was discovered."},
			{Name: "check_results", Type: proto.ColumnType_JSON, Description: "Security check results for the IP."},
			{Name: "waived_check_results", Type: proto.ColumnType_JSON, Description: "Waived security check results."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listIPs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	params := map[string]string{
		"page_size": "1000",
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

		var result ipsResponse
		if err := client.get(ctx, "/ips", params, &result); err != nil {
			return nil, fmt.Errorf("listing IPs: %w", err)
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

func getIP(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	ip := d.EqualsQualString("ip")
	if ip == "" {
		return nil, nil
	}

	params := map[string]string{
		"ip": ip,
	}

	var result IPDetails
	if err := client.get(ctx, "/ip", params, &result); err != nil {
		return nil, fmt.Errorf("getting IP %s: %w", ip, err)
	}

	return result, nil
}
