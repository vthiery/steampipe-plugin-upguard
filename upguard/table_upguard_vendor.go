package upguard

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// VendorListItem represents a vendor in the LIST /vendors response (uses mixed camelCase/snake_case).
type VendorListItem struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	PrimaryHostname  string            `json:"primary_hostname"`
	Score            int               `json:"score"`
	AutomatedScore   int               `json:"automatedScore"`
	OverallScore     int               `json:"overallScore"`
	Tier             int               `json:"tier"`
	Labels           []string          `json:"labels"`
	Monitored        bool              `json:"monitored"`
	LastAssessed     string            `json:"lastAssessed"`     // camelCase in LIST
	AssessmentStatus string            `json:"assessmentStatus"` // camelCase in LIST
	CategoryScores   CategoryScores    `json:"category_scores"`  // snake_case in LIST
	Attributes       map[string]string `json:"attributes"`
}

// Vendor represents a vendor in the GET /vendor response (uses consistent snake_case).
// This struct contains all fields available from the detailed vendor endpoint.
type Vendor struct {
	ID                   int               `json:"id"`
	Name                 string            `json:"name"`
	PrimaryHostname      string            `json:"primary_hostname"`
	Score                int               `json:"score"`
	AutomatedScore       int               `json:"automatedScore"`
	QuestionnaireScore   int               `json:"questionnaireScore"`
	OverallScore         int               `json:"overallScore"`
	IndustryAverageScore int               `json:"industry_average_score"`
	IndustryGroup        string            `json:"industry_group"`
	IndustrySector       string            `json:"industry_sector"`
	Tier                 int               `json:"tier"`
	Labels               []string          `json:"labels"`
	Portfolios           []string          `json:"portfolios"`
	Note                 string            `json:"note"`
	Monitored            bool              `json:"monitored"`
	FirstMonitored       string            `json:"first_monitored"`
	LastAssessed         string            `json:"last_assessed"`     // snake_case in GET
	ReassessmentDate     string            `json:"reassessment_date"` // only in GET
	AssessmentStatus     string            `json:"assessment_status"` // snake_case in GET
	DomainCountActive    int               `json:"domain_count_active"`
	DomainCountInactive  int               `json:"domain_count_inactive"`
	DomainCountTotal     int               `json:"domain_count_total"`
	CategoryScores       CategoryScores    `json:"categoryScores"` // camelCase in GET
	OverallRiskCounts    RiskCounts        `json:"overall_risk_counts"`
	Attributes           map[string]string `json:"attributes"`
}

// vendorsResponse is the envelope returned by GET /vendors.
type vendorsResponse struct {
	Vendors       []VendorListItem `json:"vendors"`
	NextPageToken string           `json:"next_page_token"`
	TotalResults  int              `json:"total_results"`
}

//// TABLE DEFINITION

func tableUpguardVendor() *plugin.Table {
	return &plugin.Table{
		Name:        "upguard_vendor",
		Description: "List and inspect monitored vendors in UpGuard.",
		List: &plugin.ListConfig{
			Hydrate: listVendors,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "labels", Require: plugin.Optional, Operators: []string{"="}},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AnyColumn([]string{"id", "primary_hostname"}),
			Hydrate:    getVendor,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_INT, Description: "Unique identifier for the vendor."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the vendor."},
			{Name: "primary_hostname", Type: proto.ColumnType_STRING, Description: "Primary hostname of the vendor."},
			{Name: "score", Type: proto.ColumnType_INT, Description: "Overall security score."},
			{Name: "automated_score", Type: proto.ColumnType_INT, Description: "Automated security score."},
			{Name: "questionnaire_score", Type: proto.ColumnType_INT, Description: "Questionnaire-based score."},
			{Name: "overall_score", Type: proto.ColumnType_INT, Description: "Overall combined score."},
			{Name: "industry_average_score", Type: proto.ColumnType_INT, Description: "Average score for the industry."},
			{Name: "industry_group", Type: proto.ColumnType_STRING, Description: "Industry group of the vendor."},
			{Name: "industry_sector", Type: proto.ColumnType_STRING, Description: "Industry sector of the vendor."},
			{Name: "tier", Type: proto.ColumnType_INT, Description: "Vendor tier classification."},
			{Name: "labels", Type: proto.ColumnType_JSON, Description: "Labels assigned to the vendor."},
			{Name: "portfolios", Type: proto.ColumnType_JSON, Description: "Portfolios the vendor belongs to."},
			{Name: "note", Type: proto.ColumnType_STRING, Description: "Notes about the vendor."},
			{Name: "monitored", Type: proto.ColumnType_BOOL, Description: "Whether the vendor is currently being monitored."},
		{Name: "first_monitored", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("FirstMonitored").Transform(transform.NullIfZeroValue), Description: "Date when vendor monitoring started."},
		{Name: "last_assessed", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("LastAssessed").Transform(transform.NullIfZeroValue), Description: "Date of last assessment."},
		{Name: "reassessment_date", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("ReassessmentDate").Transform(transform.NullIfZeroValue), Description: "Date when next reassessment is scheduled."},
			{Name: "assessment_status", Type: proto.ColumnType_STRING, Description: "Current assessment status (Published, Not assessed, etc.)."},
			{Name: "domain_count_active", Type: proto.ColumnType_INT, Description: "Number of active domains."},
			{Name: "domain_count_inactive", Type: proto.ColumnType_INT, Description: "Number of inactive domains."},
			{Name: "domain_count_total", Type: proto.ColumnType_INT, Description: "Total number of domains."},
			{Name: "category_scores", Type: proto.ColumnType_JSON, Description: "Breakdown of scores by category."},
			{Name: "overall_risk_counts", Type: proto.ColumnType_JSON, Description: "Count of risks by severity."},
			{Name: "attributes", Type: proto.ColumnType_JSON, Description: "Custom attributes assigned to the vendor."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listVendors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	pageToken := ""
	params := map[string]string{
		"page_size":     "1000",
		"include_risks": "false",
	}

	// Add label filter if specified
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

		var result vendorsResponse
		if err := client.get(ctx, "/vendors", params, &result); err != nil {
			return nil, fmt.Errorf("listing vendors: %w", err)
		}

		for _, listItem := range result.Vendors {
			// Convert VendorListItem to Vendor for consistent column access
			vendor := Vendor{
				ID:               listItem.ID,
				Name:             listItem.Name,
				PrimaryHostname:  listItem.PrimaryHostname,
				Score:            listItem.Score,
				AutomatedScore:   listItem.AutomatedScore,
				OverallScore:     listItem.OverallScore,
				Tier:             listItem.Tier,
				Labels:           listItem.Labels,
				Monitored:        listItem.Monitored,
				LastAssessed:     listItem.LastAssessed,
				AssessmentStatus: listItem.AssessmentStatus,
				CategoryScores:   listItem.CategoryScores,
				Attributes:       listItem.Attributes,
			}
			d.StreamListItem(ctx, vendor)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		if result.NextPageToken == "" || len(result.Vendors) == 0 {
			break
		}
		pageToken = result.NextPageToken
	}

	return nil, nil
}

func getVendor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}

	// Check if ID is provided
	if d.EqualsQuals["id"] != nil {
		params["id"] = fmt.Sprintf("%d", d.EqualsQuals["id"].GetInt64Value())
	} else if d.EqualsQuals["primary_hostname"] != nil {
		// Otherwise use hostname
		params["hostname"] = d.EqualsQuals["primary_hostname"].GetStringValue()
	} else {
		return nil, nil
	}

	var result Vendor
	if err := client.get(ctx, "/vendor", params, &result); err != nil {
		return nil, fmt.Errorf("getting vendor: %w", err)
	}

	return result, nil
}
