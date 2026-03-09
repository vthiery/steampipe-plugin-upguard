# UpGuard Plugin for Steampipe

Use SQL to query vendors, risks, domains, IPs, vulnerabilities, and breaches from [UpGuard CyberRisk](https://www.upguard.com).

## Installation

Clone and build the plugin:

```sh
git clone https://github.com/vthiery/steampipe-plugin-upguard.git
cd steampipe-plugin-upguard
make install
```

Or install from the GitHub Container Registry:

```sh
steampipe plugin install ghcr.io/vthiery/upguard
```

## Configuration

Copy the sample config:

```sh
cp config/upguard.spc ~/.steampipe/config/upguard.spc
```

Edit `~/.steampipe/config/upguard.spc`:

```hcl
connection "upguard" {
  plugin  = "upguard"

  # API key from your UpGuard CyberRisk Account Settings → API keys
  # Required API key permissions depend on the tables you query:
  # - Platform: Required for upguard_available_risk, upguard_organisation
  # - VendorRisk: Required for upguard_vendor* tables
  # - BreachRisk: Required for upguard_domain, upguard_ip, upguard_organisation_risk, upguard_vulnerability
  # - IdentityBreaches: Required for upguard_breach
  # See https://cyber-risk.upguard.com/api/docs for details
  api_key = "YOUR_API_KEY"
}
```

## API Permissions

Different tables require different API key permissions. Configure these in your UpGuard CyberRisk account settings:

| Permission | Required For |
|------------|-------------|
| **Platform** | `upguard_available_risk`, `upguard_organisation` |
| **VendorRisk** | `upguard_vendor`, `upguard_vendor_risk`, `upguard_vendor_domain`, `upguard_vendor_ip` |
| **BreachRisk** | `upguard_domain`, `upguard_ip`, `upguard_organisation_risk`, `upguard_vulnerability` |
| **IdentityBreaches** | `upguard_breach` |

Example API keys with full access would have all of these permissions: `Platform`, `VendorRisk`, `BreachRisk`, `IdentityBreaches`, `TrustExchange`, and `Admin`.



## Tables

| Table | Description |
|-------|-------------|
| [upguard_vendor](tables/upguard_vendor.md) | List and inspect monitored vendors. |
| [upguard_vendor_risk](tables/upguard_vendor_risk.md) | List active risks for a specific vendor. |
| [upguard_vendor_domain](tables/upguard_vendor_domain.md) | List domains for a specific vendor. |
| [upguard_vendor_ip](tables/upguard_vendor_ip.md) | List IP addresses for a specific vendor. |
| [upguard_domain](tables/upguard_domain.md) | List and inspect domains in your account. |
| [upguard_ip](tables/upguard_ip.md) | List and inspect IP addresses in your account. |
| [upguard_available_risk](tables/upguard_available_risk.md) | List all available risk types in the platform. |
| [upguard_organisation](tables/upguard_organisation.md) | Get information about your organization. |
| [upguard_organisation_risk](tables/upguard_organisation_risk.md) | List active risks for your organization. |
| [upguard_vulnerability](tables/upguard_vulnerability.md) | List potential vulnerabilities detected. |
| [upguard_breach](tables/upguard_breach.md) | List identity breaches detected. |

## Important Notes

### API Behavior

Some tables exhibit different behavior depending on how you query them:

- **upguard_vendor**: Querying by `id` or `primary_hostname` returns full details from the GET endpoint, including fields like `first_monitored` and `reassessment_date`. Listing vendors (e.g., by `tier`) returns summary data where some fields may be NULL.

- **upguard_domain**: Querying by `hostname` returns full details including `automated_score` and `scanned_at`. Listing domains returns only basic fields (`hostname`, `active`).

This is expected behavior based on the UpGuard API design. See [API_INCONSISTENCY.md](API_INCONSISTENCY.md) for details.

## Testing

Run the test suite to verify all tables:

```bash
./scripts/test_tables.sh
```

## Developer Documentation

- [API_INCONSISTENCY.md](API_INCONSISTENCY.md) - How the plugin handles LIST vs GET endpoint differences
- [FIELD_MAPPING_FIXES.md](FIELD_MAPPING_FIXES.md) - Summary of field mapping corrections and API investigations
