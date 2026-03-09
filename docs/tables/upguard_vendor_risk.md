# Table: upguard_vendor_risk

List active risks detected for a specific vendor in UpGuard. **Note:** You must specify `vendor_primary_hostname` in the WHERE clause.

**Required API Permission:** `VendorRisk`

## Example queries

**List all risks for a specific vendor:**

```sql
select
  risk_id,
  severity,
  category,
  detected_at,
  hostnames
from
  upguard_vendor_risk
where
  vendor_primary_hostname = 'example.com';
```

**List critical and high severity risks:**

```sql
select
  risk_id,
  severity,
  category,
  detected_at,
  hostnames
from
  upguard_vendor_risk
where
  vendor_primary_hostname = 'example.com'
  and severity in ('critical', 'high')
order by
  detected_at desc;
```

**Get risks with minimum severity level:**

```sql
select
  risk_id,
  severity,
  category,
  hostnames
from
  upguard_vendor_risk
where
  vendor_primary_hostname = 'example.com'
  and min_severity = 'high';
```

**Count risks by severity:**

```sql
select
  severity,
  count(*) as risk_count
from
  upguard_vendor_risk
where
  vendor_primary_hostname = 'example.com'
group by
  severity
order by
  risk_count desc;
```

**List risks with source details:**

```sql
select
  risk_id,
  severity,
  category,
  sources
from
  upguard_vendor_risk
where
  vendor_primary_hostname = 'example.com'
  and sources is not null;
```
