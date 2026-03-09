#!/usr/bin/env bash
# Run a smoke-test query against every UpGuard table and report results.
# Exit code is the number of failed tables (0 = all passed).

set -euo pipefail

# Colour codes (disabled if not a terminal)
if [ -t 1 ]; then
  GREEN="\033[0;32m"
  RED="\033[0;31m"
  YELLOW="\033[0;33m"
  RESET="\033[0m"
else
  GREEN="" RED="" YELLOW="" RESET=""
fi

PASS=0
FAIL=0
SKIP=0

run_test() {
  local table="$1"
  local query="$2"
  printf "  %-40s" "$table"
  local output
  if output=$(steampipe query "$query" --output json 2>&1); then
    printf "${GREEN}PASS${RESET}\n"
    ((PASS++)) || true
  else
    # Check for authentication or permission errors
    if echo "$output" | grep -qE "(403|401|Forbidden|Unauthorized)"; then
      printf "${YELLOW}SKIP${RESET} (authentication/permission issue)\n"
      ((SKIP++)) || true
    else
      printf "${RED}FAIL${RESET}\n"
      echo "$output" | sed 's/^/    /'
      ((FAIL++)) || true
    fi
  fi
}

echo ""
echo "UpGuard Steampipe plugin — table smoke tests"
echo "================================================="
echo ""

# Get a vendor hostname for vendor-specific table tests
echo "Fetching vendor for vendor-specific tests..."
VENDOR_HOSTNAME=$(steampipe query "select primary_hostname from upguard_vendor limit 1" --output json 2>/dev/null | grep -o '"primary_hostname"[[:space:]]*:[[:space:]]*"[^"]*"' | head -1 | sed 's/.*"primary_hostname"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/' || echo "")

if [ -z "$VENDOR_HOSTNAME" ]; then
  echo "  ${YELLOW}Warning:${RESET} Could not fetch vendor hostname. Vendor-specific tests will be skipped."
else
  echo "  Found vendor: $VENDOR_HOSTNAME"
fi

echo ""
echo "Testing tables..."
echo ""

# Test organisation table
run_test "upguard_organisation" \
  "select name, primary_hostname, automated_score from upguard_organisation"

# Test available risks
run_test "upguard_available_risk" \
  "select risk_type, severity, category from upguard_available_risk limit 5"

# Test organisation risks
run_test "upguard_organisation_risk" \
  "select risk_id, severity, category, detected_at from upguard_organisation_risk limit 5"

# Test vendors
run_test "upguard_vendor" \
  "select name, monitored, assessment_status from upguard_vendor limit 5"

# Test domains
run_test "upguard_domain" \
  "select hostname, automated_score, active, scanned_at from upguard_domain limit 5"

# Test IPs
run_test "upguard_ip" \
  "select ip, country, automated_score, asn from upguard_ip limit 5"

# Test vulnerabilities (with new field structure)
run_test "upguard_vulnerability" \
  "select cve, severity, cvss, created_at from upguard_vulnerability limit 5"

# Test breaches
run_test "upguard_breach" \
  "select title, date_occurred, breach_type from upguard_breach limit 5"

# Vendor-specific tables (require vendor_primary_hostname)
if [ -n "$VENDOR_HOSTNAME" ]; then
  run_test "upguard_vendor_risk" \
    "select vendor_primary_hostname, risk_id, severity, detected_at from upguard_vendor_risk where vendor_primary_hostname = '$VENDOR_HOSTNAME' limit 5"
  
  run_test "upguard_vendor_domain" \
    "select vendor_primary_hostname, hostname, automated_score, active from upguard_vendor_domain where vendor_primary_hostname = '$VENDOR_HOSTNAME' limit 5"
  
  run_test "upguard_vendor_ip" \
    "select vendor_primary_hostname, ip, country, automated_score from upguard_vendor_ip where vendor_primary_hostname = '$VENDOR_HOSTNAME' limit 5"
else
  printf "  %-40s${YELLOW}SKIP${RESET} (no vendor hostname)\n" "upguard_vendor_risk"
  printf "  %-40s${YELLOW}SKIP${RESET} (no vendor hostname)\n" "upguard_vendor_domain"
  printf "  %-40s${YELLOW}SKIP${RESET} (no vendor hostname)\n" "upguard_vendor_ip"
  SKIP=$((SKIP + 3))
fi

echo ""
echo "-------------------------------------------------"
printf "Results: ${GREEN}%d passed${RESET}  ${RED}%d failed${RESET}  ${YELLOW}%d skipped${RESET}\n" "$PASS" "$FAIL" "$SKIP"
echo ""

exit "$FAIL"
