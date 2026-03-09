# Table: upguard_ip

List and inspect IP addresses in your UpGuard account.

**Required API Permission:** `BreachRisk`

## Example queries

**List all IPs:**

```sql
select
  ip,
  hostname,
  asn,
  asn_description
from
  upguard_ip;
```

**Get detailed information for a specific IP:**

```sql
select
  ip,
  hostname,
  asn,
  asn_description,
  services,
  check_results
from
  upguard_ip
where
  ip = '192.0.2.1';
```

**List IPs grouped by ASN:**

```sql
select
  asn,
  asn_description,
  count(*) as ip_count
from
  upguard_ip
group by
  asn, asn_description
order by
  ip_count desc;
```

**List IPs with their associated services:**

```sql
select
  ip,
  hostname,
  jsonb_array_length(services) as service_count
from
  upguard_ip
where
  services is not null
order by
  service_count desc;
```
