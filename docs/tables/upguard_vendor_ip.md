# Table: upguard_vendor_ip

List IP addresses for a specific vendor. **Note:** You must specify `vendor_primary_hostname` in the WHERE clause.

**Required API Permission:** `VendorRisk`

## Example queries

**List all IPs for a specific vendor:**

```sql
select
  ip,
  owner,
  asn,
  as_name
from
  upguard_vendor_ip
where
  vendor_primary_hostname = 'example.com';
```

**Get detailed IP information:**

```sql
select
  ip,
  owner,
  asn,
  as_name,
  services,
  check_results
from
  upguard_vendor_ip
where
  vendor_primary_hostname = 'example.com'
  and ip = '192.0.2.1';
```

**List IPs with their ASN details:**

```sql
select
  ip,
  owner,
  asn,
  as_name
from
  upguard_vendor_ip
where
  vendor_primary_hostname = 'example.com'
order by
  asn;
```
