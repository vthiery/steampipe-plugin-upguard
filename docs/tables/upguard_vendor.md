# Table: upguard_vendor

List and inspect monitored vendors in your UpGuard CyberRisk account.

**Required API Permission:** `VendorRisk`

## Example queries

**List all monitored vendors:**

```sql
select
  id,
  name,
  primary_hostname,
  score,
  tier
from
  upguard_vendor;
```

**List vendors with critical risk counts:**

```sql
select
  name,
  primary_hostname,
  score,
  overall_risk_counts
from
  upguard_vendor
where
  (overall_risk_counts->>'critical')::int > 0
order by
  (overall_risk_counts->>'critical')::int desc;
```

**List vendors by tier:**

```sql
select
  name,
  primary_hostname,
  score,
  tier,
  labels
from
  upguard_vendor
where
  tier = 1
order by
  score desc;
```

**Get vendor details by hostname:**

```sql
select
  name,
  score,
  automated_score,
  questionnaire_score,
  industry_group,
  category_scores
from
  upguard_vendor
where
  primary_hostname = 'example.com';
```

**List vendors with specific labels:**

```sql
select
  name,
  primary_hostname,
  score,
  labels
from
  upguard_vendor
where
  labels @> '["critical"]';
```

## Important Notes

### LIST vs GET Behavior

This table exhibits different behavior depending on the query:

- **Querying by id or primary_hostname** (e.g., `WHERE id = 12345` or `WHERE primary_hostname = 'example.com'`): Returns full details from the GET `/vendor` endpoint, including `first_monitored`, `reassessment_date`, `domain_count_*`, `industry_*`, `portfolios`, and `note`. The `monitored` field will show `false` (not returned by GET endpoint).

- **Listing vendors** (e.g., `WHERE tier = 1` or no WHERE clause): Returns summary data from the LIST `/vendors` endpoint. Fields like `first_monitored`, `reassessment_date`, `domain_count_*`, `industry_*`, `portfolios`, and `note` will be NULL (not included in LIST response). The `monitored` field will show `true` for monitored vendors.

This is expected behavior based on the UpGuard API design. For full vendor details, query by specific id or hostname.

