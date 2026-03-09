# Table: upguard_breach

List identity breaches detected by UpGuard.

**Required API Permission:** `IdentityBreaches`

## Example queries

**List all breaches:**

```sql
select
  name,
  type,
  breach_date,
  exposed_record_count
from
  upguard_breach
order by
  breach_date desc;
```

**List breaches with exposed data classes:**

```sql
select
  name,
  breach_date,
  exposed_record_count,
  exposed_data_classes
from
  upguard_breach
order by
  exposed_record_count desc;
```

**Get detailed information for a specific breach:**

```sql
select
  name,
  type,
  breach_date,
  exposed_record_count,
  exposed_data_classes,
  description
from
  upguard_breach
where
  name = 'Example Breach';
```

**Count breaches by year:**

```sql
select
  extract(year from breach_date) as breach_year,
  count(*) as breach_count
from
  upguard_breach
group by
  breach_year
order by
  breach_year desc;
```

**List breaches with personally identifiable information:**

```sql
select
  name,
  breach_date,
  exposed_record_count,
  exposed_data_classes
from
  upguard_breach
where
  exposed_data_classes::text like '%email%'
  or exposed_data_classes::text like '%password%'
order by
  breach_date desc;
```
