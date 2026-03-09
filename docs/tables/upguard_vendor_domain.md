# Table: upguard_vendor_domain

List domains for a specific vendor. **Note:** You must specify `vendor_primary_hostname` in the WHERE clause.

**Required API Permission:** `VendorRisk`

## Example queries

**List all domains for a specific vendor:**

```sql
select
  hostname,
  active,
  labels
from
  upguard_vendor_domain
where
  vendor_primary_hostname = 'example.com';
```

**List active domains for a vendor:**

```sql
select
  hostname,
  active,
  labels
from
  upguard_vendor_domain
where
  vendor_primary_hostname = 'example.com'
  and active = true;
```

**Get detailed domain information:**

```sql
select
  hostname,
  active,
  automated_score,
  scanned_at,
  a_records,
  check_results
from
  upguard_vendor_domain
where
  vendor_primary_hostname = 'example.com'
  and hostname = 'subdomain.example.com';
```
