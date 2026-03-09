# UpGuard Plugin for Steampipe

Use SQL to query vendors, risks, domains, IPs, vulnerabilities, and breaches from [UpGuard CyberRisk](https://www.upguard.com).

- **[Get started →](https://github.com/vthiery/steampipe-plugin-upguard)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-upguard/issues)

## Quick start

### Install

Install the plugin with [Steampipe](https://steampipe.io):

```sh
steampipe plugin install ghcr.io/vthiery/upguard
```

### Configure

Copy the sample config and set your API key:

```sh
cp config/upguard.spc ~/.steampipe/config/upguard.spc
```

Edit `~/.steampipe/config/upguard.spc`:

```hcl
connection "upguard" {
  plugin  = "ghcr.io/vthiery/upguard"

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

### API Permissions

Different tables require different API key permissions:

| Permission | Required For |
|------------|-------------|
| **Platform** | `upguard_available_risk`, `upguard_organisation` |
| **VendorRisk** | `upguard_vendor`, `upguard_vendor_risk`, `upguard_vendor_domain`, `upguard_vendor_ip` |
| **BreachRisk** | `upguard_domain`, `upguard_ip`, `upguard_organisation_risk`, `upguard_vulnerability` |
| **IdentityBreaches** | `upguard_breach` |

You can configure API key permissions in your UpGuard CyberRisk account settings.

### Run a query

```shell
steampipe query
```

```sql
-- List all monitored vendors with their security scores
select
  name,
  primary_hostname,
  score,
  tier,
  industry_group
from
  upguard_vendor
order by
  score desc;
```

## Tables

| Table | Description |
|-------|-------------|
| [upguard_vendor](docs/tables/upguard_vendor.md) | List and inspect monitored vendors. |
| [upguard_vendor_risk](docs/tables/upguard_vendor_risk.md) | List active risks for a specific vendor. |
| [upguard_vendor_domain](docs/tables/upguard_vendor_domain.md) | List domains for a specific vendor. |
| [upguard_vendor_ip](docs/tables/upguard_vendor_ip.md) | List IP addresses for a specific vendor. |
| [upguard_domain](docs/tables/upguard_domain.md) | List and inspect domains in your account. |
| [upguard_ip](docs/tables/upguard_ip.md) | List and inspect IP addresses in your account. |
| [upguard_available_risk](docs/tables/upguard_available_risk.md) | List all available risk types in the platform. |
| [upguard_organisation](docs/tables/upguard_organisation.md) | Get information about your organization. |
| [upguard_organisation_risk](docs/tables/upguard_organisation_risk.md) | List active risks for your organization. |
| [upguard_vulnerability](docs/tables/upguard_vulnerability.md) | List potential vulnerabilities detected. |
| [upguard_breach](docs/tables/upguard_breach.md) | List identity breaches detected. |

## Development

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
make install
```

Configure the plugin:

```sh
cp config/upguard.spc ~/.steampipe/config/upguard.spc
vi ~/.steampipe/config/upguard.spc
```

## Testing

Run a smoke query against every table:

```sh
make test
```

The test script ([scripts/test_tables.sh](scripts/test_tables.sh)) builds the plugin, queries each table, and reports pass/fail/skip (scope-restricted tables are skipped rather than failed).

### Further reading

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [UpGuard API reference](https://cyber-risk.upguard.com/api/docs)
