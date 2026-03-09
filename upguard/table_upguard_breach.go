package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// BreachedIdentity represents a breached identity.
type BreachedIdentity struct {
	Name           string `json:"name"`
	Domain         string `json:"domain"`
	Email          string `json:"email"`
	NumBreaches    int    `json:"num_breaches"`
	DateLastBreach string `json:"date_last_breach"`
}

// Breach represents a data breach.
type Breach struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Title              string   `json:"title"`
	Domain             string   `json:"domain"`
	BreachType         string   `json:"breach_type"`
	DateOccurred       string   `json:"date_occurred"`
	DatePublished      string   `json:"date_published"`
	Description        string   `json:"description"`
	TotalExposures     int      `json:"total_exposures"`
	ExposedDataClasses []string `json:"exposed_data_classes"`
	AssigneeUserEmail  string   `json:"assignee_user_email"`
}

// breachesResponse is the envelope returned by GET /breaches.
type breachesResponse struct {
	BreachedIdentities []BreachedIdentity `json:"breached_identities"`
	Breaches           []Breach           `json:"breaches"`
	NextPageToken      string             `json:"next_page_token"`
	TotalResults       int                `json:"total_results"`
}

// breachResponse is the envelope returned by GET /breach.
type breachResponse struct {
	Breach   Breach                   `json:"breach"`
	Comments []map[string]interface{} `json:"comments"`
}

//// TABLE DEFINITION

func tableUpguardBreach() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_breach",
		Description: "List identity breaches detected in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listBreaches,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "breach_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getBreach,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the breach."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the breach."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "Title of the breach."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain affected by the breach."},
			{Name: "breach_type", Type: proto.ColumnType_STRING, Description: "Type of breach (Company, Website, etc)."},
			{Name: "date_occurred", Type: proto.ColumnType_TIMESTAMP, Description: "When the breach occurred."},
			{Name: "date_published", Type: proto.ColumnType_TIMESTAMP, Description: "When the breach was published."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the breach."},
			{Name: "total_exposures", Type: proto.ColumnType_INT, Description: "Total number of exposures."},
			{Name: "exposed_data_classes", Type: proto.ColumnType_JSON, Description: "Types of data exposed in the breach."},
			{Name: "assignee_user_email", Type: proto.ColumnType_STRING, Description: "Email of the user assigned to the breach."},
			{Name: "breach_id", Type: proto.ColumnType_STRING, Transform: transform.FromConstant(""), Description: "Breach ID filter for list queries."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listBreaches(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	params := map[string]string{
		"page_size": "1000",
	}

	// Add breach_id filter if specified
	if d.EqualsQuals["breach_id"] != nil {
		params["breach_id"] = d.EqualsQualString("breach_id")
	}

	for {
		if pageToken != "" {
			params["page_token"] = pageToken
		}

		var result breachesResponse
		if err := client.get(ctx, "/breaches", params, &result); err != nil {
			return nil, fmt.Errorf("listing breaches: %w", err)
		}

		for _, breach := range result.Breaches {
			d.StreamListItem(ctx, breach)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextPageToken == "" || len(result.Breaches) == 0 {
			break
		}
		pageToken = result.NextPageToken
	}

	return nil, nil
}

func getBreach(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	id := d.EqualsQuals["id"].GetInt64Value()
	if id == 0 {
		return nil, nil
	}

	params := map[string]string{
		"id": fmt.Sprintf("%d", id),
	}

	var result breachResponse
	if err := client.get(ctx, "/breach", params, &result); err != nil {
		return nil, fmt.Errorf("getting breach %d: %w", id, err)
	}

	return result.Breach, nil
}
