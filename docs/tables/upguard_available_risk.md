# Table: upguard_available_risk

List all available risk types in the UpGuard platform with their descriptions and remediation guidance.

**Required API Permission:** `Platform`

## Example queries

**List all available risks:**

```sql
select
  id,
  risk,
  severity,
  category,
  group
from
  upguard_available_risk
order by
  severity, category;
```

**Get details for a specific risk:**

```sql
select
  id,
  risk,
  finding,
  risk_details,
  remediation,
  severity
from
  upguard_available_risk
where
  id = 'domain_expired';
```

**List critical risks by category:**

```sql
select
  category,
  count(*) as risk_count
from
  upguard_available_risk
where
  severity = 'critical'
group by
  category
order by
  risk_count desc;
```

**List risks with remediation guidance:**

```sql
select
  id,
  risk,
  severity,
  remediation
from
  upguard_available_risk
where
  severity in ('high', 'critical')
  and remediation is not null
order by
  severity, risk;
```

**Find risks related to SSL/TLS:**

```sql
select
  id,
  risk,
  severity,
  category,
  description
from
  upguard_available_risk
where
  description ilike '%ssl%'
  or description ilike '%tls%'
  or risk ilike '%ssl%'
  or risk ilike '%tls%';
```
