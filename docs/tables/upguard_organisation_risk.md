# Table: upguard_organisation_risk

List active risks for your organization in UpGuard.

**Required API Permission:** `BreachRisk`

## Example queries

**List all organization risks:**

```sql
select
  risk_id,
  hostname,
  severity,
  status
from
  upguard_organisation_risk
order by
  severity;
```

**List critical and high severity risks:**

```sql
select
  risk_id,
  hostname,
  severity,
  status,
  opened_at
from
  upguard_organisation_risk
where
  severity in ('critical', 'high')
order by
  opened_at desc;
```

**Count risks by severity:**

```sql
select
  severity,
  count(*) as risk_count
from
  upguard_organisation_risk
group by
  severity
order by
  risk_count desc;
```

**Get detailed risk information:**

```sql
select
  risk_id,
  hostname,
  severity,
  status,
  sources,
  metadata
from
  upguard_organisation_risk
where
  risk_id = 'domain_expired';
```
