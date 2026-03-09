#!/bin/bash

# Test script to verify vendor field population in LIST vs GET queries
# This demonstrates the UpGuard API inconsistency resolution

echo "========================================="
echo "Testing UpGuard Vendor Field Population"
echo "========================================="
echo ""

echo "1. LIST Query (tier=1) - Shows fields available from LIST endpoint"
echo "   Expected: monitored, assessment_status, last_assessed populated"
echo "   Expected: first_monitored, reassessment_date = NULL (not in LIST endpoint)"
echo ""
steampipe query "select name, monitored, assessment_status, last_assessed, first_monitored, reassessment_date from upguard_vendor where tier = 1 limit 3"

echo ""
echo "2. GET Query (by hostname) - Shows ALL fields from detailed endpoint"
echo "   Expected: first_monitored, reassessment_date populated"
echo "   Note: monitored=false (not returned by GET endpoint, shows zero value)"
echo ""
steampipe query "select name, first_monitored, last_assessed, reassessment_date, assessment_status from upguard_vendor where primary_hostname = 'google.com'"

echo ""
echo "3. GET Query (by id) - Shows ALL fields from detailed endpoint"
echo "   Expected: first_monitored, reassessment_date populated"
echo "   Note: monitored=false (not returned by GET endpoint, shows zero value)"
echo ""
steampipe query "select name, first_monitored, last_assessed, reassessment_date, assessment_status from upguard_vendor where id = 5360402938462208"

echo ""
echo "========================================="
echo "Summary"
echo "========================================="
echo "The plugin handles UpGuard API inconsistency by:"
echo "1. Using VendorListItem struct for LIST /vendors (camelCase: assessmentStatus, lastAssessed)"
echo "2. Using Vendor struct for GET /vendor (snake_case: assessment_status, last_assessed)"
echo "3. Converting between structs in listVendors() function"
echo ""
echo "Result:"
echo "- LIST queries show fields available in LIST endpoint (monitored=true, dates, status)"
echo "- GET queries (by id/hostname) show ALL fields including first_monitored, reassessment_date"
echo "- No data loss, proper handling of API inconsistency"
echo ""
