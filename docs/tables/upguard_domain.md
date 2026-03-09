# Table: upguard_domain

List and inspect domains in your UpGuard account.

**Required API Permission:** `BreachRisk`

## Example queries

**List all active domains:**

```sql
select
  hostname,
  active,
  automated_score,
  scanned_at,
  labels
from
  upguard_domain
where
  active = true
order by
  automated_score;
```

**Get detailed information for a specific domain:**

```sql
select
  hostname,
  automated_score,
  scanned_at,
  a_records,
  check_results
from
  upguard_domain
where
  hostname = 'example.com';
```

**List domains with low security scores:**

```sql
select
  hostname,
  automated_score,
  scanned_at,
  labels
from
  upguard_domain
where
  automated_score < 700
order by
  automated_score;
```

**List domains by label:**

```sql
select
  hostname,
  automated_score,
  active,
  labels
from
  upguard_domain
where
  labels @> '["production"]'
order by
  hostname;
```

**Count domains by active status:**

```sql
select
  active,
  count(*) as domain_count
from
  upguard_domain
group by
  active;
```

## Important Notes

### LIST vs GET Behavior

This table exhibits different behavior depending on the query:

- **Querying by hostname** (e.g., `WHERE hostname = 'example.com'`): Returns full details from the GET `/domain` endpoint, including `automated_score`, `scanned_at`, `a_records`, `check_results`, and `labels`.

- **Listing domains** (e.g., `WHERE active = true` or no WHERE clause): Returns basic information from the LIST `/domains` endpoint, including only `hostname` and `active`. Fields like `automated_score` and `scanned_at` will be NULL.

This is expected behavior based on the UpGuard API design. For full domain details, query by specific hostname.
