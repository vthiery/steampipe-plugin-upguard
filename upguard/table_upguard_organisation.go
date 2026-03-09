package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TYPES

// Organisation represents the current organization.
type Organisation struct {
	ID              int            `json:"id"`
	Name            string         `json:"name"`
	PrimaryHostname string         `json:"primary_hostname"`
	AutomatedScore  int            `json:"automatedScore"`
	OverallScore    int            `json:"overallScore"`
	CategoryScores  CategoryScores `json:"categoryScores"`
}

//// TABLE DEFINITION

func tableUpguardOrganisation() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_organisation",
		Description: "Get information about your UpGuard organization.",
		List: &plugin.ListConfig{
			Hydrate: listOrganisation,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the organization."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the organization."},
			{Name: "primary_hostname", Type: proto.ColumnType_STRING, Description: "Primary hostname of the organization."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score."},
			{Name: "overall_score", Type: proto.ColumnType_INT, Description: "Overall security score."},
			{Name: "category_scores", Type: proto.ColumnType_JSON, Description: "Breakdown of scores by category."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listOrganisation(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var result Organisation
	if err := client.get(ctx, "/organisation", nil, &result); err != nil {
		return nil, fmt.Errorf("getting organisation: %w", err)
	}

	d.StreamListItem(ctx, result)

	return nil, nil
}
