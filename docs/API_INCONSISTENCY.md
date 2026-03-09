# UpGuard API Inconsistency Resolution

## Overview

The UpGuard API returns different field name formats and different field sets between LIST and GET endpoints. This is a consistent pattern across multiple resource types (vendors, domains, etc.). This document explains how the plugin handles these inconsistencies.

## Affected Tables

The following tables use separate structs to handle LIST vs GET inconsistencies:

1. **upguard_vendor** - Different field names and field sets
2. **upguard_domain** - Different field sets (LIST returns minimal data)

## Problem Patterns

### Pattern 1: Different Field Names (Vendor Table)

The UpGuard API returns different field name formats between LIST and GET endpoints:

- **LIST endpoint** (`/vendors`):
  - Uses camelCase for some fields: `assessmentStatus`, `lastAssessed`
  - Uses snake_case for others: `category_scores`
  - Returns 13 fields including `monitored` (bool)
  - Does NOT return: `first_monitored`, `reassessment_date`, `domain_count_*`, `industry_*`, etc.

- **GET endpoint** (`/vendor`):
  - Uses snake_case for most fields: `assessment_status`, `last_assessed`
  - Uses camelCase for category scores: `categoryScores`
  - Returns 27+ fields including all LIST fields PLUS additional detail fields
  - Does NOT return: `monitored` field

## Solution

The plugin uses separate Go structs to handle each endpoint's response format, then converts between them in the List hydrate function. This ensures all columns have consistent data types and field mappings regardless of which endpoint was called.

### Example 1: Vendor Table (Different Field Names)

### 1. VendorListItem Struct
```go
type VendorListItem struct {
    ID               int            `json:"id"`
    Name             string         `json:"name"`
    // ... other fields ...
    LastAssessed     string         `json:"lastAssessed"`     // camelCase
    AssessmentStatus string         `json:"assessmentStatus"` // camelCase
    CategoryScores   CategoryScores `json:"category_scores"`  // snake_case
    Monitored        bool           `json:"monitored"`        // only in LIST
}
```

### 2. Vendor Struct
```go
type Vendor struct {
    ID                   int            `json:"id"`
    Name                 string         `json:"name"`
    // ... other fields ...
    FirstMonitored       string         `json:"first_monitored"`   // only in GET
    LastAssessed         string         `json:"last_assessed"`     // snake_case
    ReassessmentDate     string         `json:"reassessment_date"` // only in GET
    AssessmentStatus     string         `json:"assessment_status"` // snake_case
    CategoryScores       CategoryScores `json:"categoryScores"`    // camelCase
    Monitored            bool           `json:"monitored"`         // NOT in GET (zero value)
}
```

### 3. Conversion in listVendors()
```go
for _, listItem := range result.Vendors {
    // Convert VendorListItem to Vendor for consistent column access
    vendor := Vendor{
        ID:               listItem.ID,
        Name:             listItem.Name,
        // ... copy all available fields ...
        LastAssessed:     listItem.LastAssessed,
        AssessmentStatus: listItem.AssessmentStatus,
        // Fields not in LIST remain zero values:
        // FirstMonitored, ReassessmentDate, etc.
    }
    d.StreamListItem(ctx, vendor)
}
```

### Example 2: Domain Table (Different Field Sets)

**LIST endpoint** (`/domains`):
- Returns only: `hostname`, `active`, `primary_domain`
- No scoring, scanning, or check data

**GET endpoint** (`/domain`):
- Returns all LIST fields PLUS: `automated_score`, `scanned_at`, `a_records`, `check_results`, `labels`, etc.

**Solution**:
```go
// DomainListItem for LIST endpoint
type DomainListItem struct {
    Hostname      string `json:"hostname"`
    Active        bool   `json:"active"`
    PrimaryDomain bool   `json:"primary_domain"`
}

// Domain for GET endpoint (full details)
type Domain struct {
    Hostname           string        `json:"hostname"`
    Active             bool          `json:"active"`
    AutomatedScore     int           `json:"automated_score"`
    ScannedAt          string        `json:"scanned_at"`
    ARecords           []string      `json:"a_records"`
    Labels             []string      `json:"labels"`
    CheckResults       []CheckResult `json:"check_results"`
    WaivedCheckResults []CheckResult `json:"waived_check_results"`
}

// Conversion in listDomains()
for _, listItem := range result.Domains {
    domain := Domain{
        Hostname: listItem.Hostname,
        Active:   listItem.Active,
        // AutomatedScore, ScannedAt, etc. remain zero values
    }
    d.StreamListItem(ctx, domain)
}
```

## Results

### Vendor Table

#### LIST Queries (e.g., `WHERE tier = 1`)
- Populated fields: `monitored`, `assessment_status`, `last_assessed`, etc.
- NULL fields: `first_monitored`, `reassessment_date` (not in LIST endpoint)

#### GET Queries (e.g., `WHERE id = X` or `WHERE primary_hostname = 'example.com'`)
- Populated fields: All fields including `first_monitored`, `reassessment_date`, `domain_count_*`, etc.
- Zero value fields: `monitored` shows `false` (not returned by GET endpoint)

### Domain Table

#### LIST Queries (e.g., `WHERE active = true`)
- Populated fields: `hostname`, `active` only
- NULL fields: `automated_score`, `scanned_at`, `a_records`, etc. (not in LIST endpoint)

#### GET Queries (e.g., `WHERE hostname = 'example.com'`)
- Populated fields: All fields including `automated_score`, `scanned_at`, `check_results`, etc.

## Testing

Run the test scripts to verify:
```bash
# Test vendor-specific LIST vs GET behavior
./scripts/test_vendor_fields.sh

# Test all tables including fixed field mappings
./scripts/test_tables.sh
```

After making changes to table files, rebuild and restart Steampipe:
```bash
make install
steampipe service restart --force
```

## Related Issues

See [FIELD_MAPPING_FIXES.md](FIELD_MAPPING_FIXES.md) for details on other field mapping issues that were fixed, including:
- upguard_vulnerability - Incorrect field structure
- upguard_organisation - Non-existent fields
- CheckResult struct - Wrong field types

## Summary

- **Problem**: UpGuard API returns different data structures between LIST and GET endpoints
- **Solution**: Use separate structs per endpoint and convert in hydrate functions
- **Result**: No data loss, all fields accessible, consistent column definitions
- **Pattern**: Applies to vendor and domain tables; other tables may use single endpoint
- **Expected Behavior**: Fields show NULL when not available from the queried endpoint
