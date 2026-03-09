// Package upguard provides a Steampipe plugin for querying UpGuard CyberRisk resources using SQL.
package upguard

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin returns the definition of the UpGuard Steampipe plugin.
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-upguard",
		DefaultTransform: transform.FromGo().NullIfZero(),
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"upguard_vendor":            tableUpguardVendor(),
			"upguard_vendor_risk":       tableUpguardVendorRisk(),
			"upguard_vendor_domain":     tableUpguardVendorDomain(),
			"upguard_vendor_ip":         tableUpguardVendorIP(),
			"upguard_domain":            tableUpguardDomain(),
			"upguard_ip":                tableUpguardIP(),
			"upguard_available_risk":    tableUpguardAvailableRisk(),
			"upguard_organisation":      tableUpguardOrganisation(),
			"upguard_organisation_risk": tableUpguardOrganisationRisk(),
			"upguard_vulnerability":     tableUpguardVulnerability(),
			"upguard_breach":            tableUpguardBreach(),
		},
	}
	return p
}
