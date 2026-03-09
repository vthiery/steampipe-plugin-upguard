# Table: upguard_organisation

Get information about your UpGuard organization. This table returns a single row with organization details.

**Required API Permission:** `Platform`

## Example queries

**Get organization details:**

```sql
select
  primary_hostname,
  name,
  score,
  overall_risk_counts
from
  upguard_organisation;
```

**Get organization risk breakdown:**

```sql
select
  primary_hostname,
  name,
  overall_risk_counts->>'critical' as critical_risks,
  overall_risk_counts->>'high' as high_risks,
  overall_risk_counts->>'medium' as medium_risks,
  overall_risk_counts->>'low' as low_risks
from
  upguard_organisation;
```

**Get organization score and category scores:**

```sql
select
  primary_hostname,
  name,
  score,
  category_scores
from
  upguard_organisation;
```
