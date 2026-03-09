# Field Mapping Fixes Summary

## Issues Found and Fixed

### 1. upguard_domain Table ✅ FIXED
**Problem**: LIST `/domains` endpoint returns different fields than GET `/domain`
- LIST returns: `hostname`, `active`, `primary_domain` only
- GET returns: All LIST fields PLUS `automated_score`, `scanned_at`, `a_records`, `check_results`, etc.

**Solution**: Implemented separate structs similar to vendor:
- `DomainListItem` for LIST endpoint (minimal fields)
- `Domain` for GET endpoint (full details)
- Conversion logic in `listDomains()` function

**Test Results**:
```sql
select hostname, automated_score, scanned_at from upguard_domain where hostname = 'camunda.com'
-- Result: camunda.com | 789 | 2026-03-09T00:36:49Z ✅
```

### 2. upguard_vulnerability Table ✅ FIXED
**Problem**: Struct fields completely mismatched with actual API response
- Expected: `CVE` (map), `CVSS` (float), `Severity` (string), `DetectedAt`, `Hostnames` (array), `IPs` (array)
- Actual API returns: `hostname` (string), `ip_addresses` (array), `cve` (object with id/description/severity/epss), `created_at`, `verified`, `known_exploited_vulnerability`, `cpes`

**Solution**: Rewrote entire struct to match API:
- Created `CVEInfo` struct for nested CVE object
- Updated `Vulnerability` struct with correct field names and types
- Added transform to convert CVSS score to severity level
- Updated all column definitions

**Test Results**:
```sql
select cve, cvss, severity, hostname, created_at from upguard_vulnerability limit 3
-- Result: CVE-2020-7656 | 6.1 | medium | page.camunda.com | 2023-08-02T12:23:14Z ✅
```

### 3. upguard_organisation Table ✅ FIXED
**Problem**: `overallScore` field defined in struct but not returned by API
- API returns: `id`, `name`, `primary_hostname`, `automatedScore`, `categoryScores`
- Struct expected: All above PLUS `overallScore` (which doesn't exist)

**Solution**: Removed `OverallScore` field from struct and `overall_score` column from table

**Test Results**:
```sql
select name, automated_score, category_scores from upguard_organisation
-- Result: Camunda | 899 | {"attackSurface":950, ...} ✅
```

### 4. CheckResult Struct ✅ FIXED
**Problem**: CheckResult struct used by domain table had incorrect field types and names
- Expected: `risk_id` (string), `severity` (string), `detected_at`, `last_scanned_at`
- Actual API: `id`, `riskType`, `severity` (int), `severityName` (string), `checked_at`, `title`, `description`, etc.

**Solution**: Rewrote CheckResult struct to match actual API response with all fields and correct types

## Summary of Changes

### Files Modified:
1. [upguard/table_upguard_domain.go](upguard/table_upguard_domain.go)
   - Added `DomainListItem` struct for LIST endpoint
   - Renamed `DomainDetails` to `Domain` for GET endpoint
   - Fixed `CheckResult` struct with correct fields
   - Added conversion logic in `listDomains()`

2. [upguard/table_upguard_vulnerability.go](upguard/table_upguard_vulnerability.go)
   - Added `CVEInfo` struct for nested CVE data
   - Completely rewrote `Vulnerability` struct
   - Updated all column definitions with transforms
   - Added `cvssToSeverityTransform()` function

3. [upguard/table_upguard_organisation.go](upguard/table_upguard_organisation.go)
   - Removed `OverallScore` field from struct
   - Removed `overall_score` column from table

4. [upguard/table_upguard_vendor_domain.go](upguard/table_upguard_vendor_domain.go)
   - Updated reference from `DomainDetails` to `Domain`

## Testing Results

All tables now correctly handle API responses:

| Table | Status | Notes |
|-------|--------|-------|
| upguard_vendor | ✅ PASS | Already fixed (LIST vs GET inconsistency) |
| upguard_domain | ✅ PASS | Fixed LIST vs GET inconsistency |
| upguard_vulnerability | ✅ PASS | Fixed all field mappings |
| upguard_organisation | ✅ PASS | Removed non-existent field |
| upguard_vendor_domain | ✅ PASS | Updated to use new Domain struct |
| upguard_vendor_ip | ✅ PASS | No changes needed |
| upguard_vendor_risk | ✅ PASS | No changes needed |
| upguard_ip | ✅ PASS | No changes needed |
| upguard_organisation_risk | ✅ PASS | No changes needed |
| upguard_breach | ✅ PASS | No changes needed |
| upguard_available_risk | ✅ PASS | No changes needed |

## API Inconsistencies Documented

### UpGuard API Patterns:
1. **LIST vs GET Endpoints**: Many UpGuard API endpoints return different fields between LIST and GET:
   - LIST endpoints return minimal fields for pagination efficiency
   - GET endpoints return full details
   - Solution: Use separate structs and convert in hydrate functions

2. **Field Naming**: UpGuard API uses inconsistent naming:
   - Mix of camelCase (`automatedScore`) and snake_case (`primary_hostname`)
   - Different naming between LIST and GET (vendor: `assessmentStatus` in LIST vs `assessment_status` in GET)

3. **Nested Objects**: Some fields like `cve` in vulnerabilities are complex objects requiring separate structs

## Recommendations

1. ✅ All field mappings verified against actual API responses
2. ✅ Separate structs used for LIST vs GET where needed
3. ✅ Transform functions added for data type conversions
4. ✅ Documentation updated to reflect actual API behavior
5. 📝 Consider creating integration tests that verify API response structure hasn't changed

## Related Documentation
- [API_INCONSISTENCY.md](API_INCONSISTENCY.md) - Detailed explanation of LIST vs GET patterns
- [scripts/test_vendor_fields.sh](../scripts/test_vendor_fields.sh) - Test script for vendor field population
- [scripts/test_tables.sh](../scripts/test_tables.sh) - Comprehensive test script for all tables
