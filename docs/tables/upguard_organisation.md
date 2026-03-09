# Table: upguard_organisation

Get information about your UpGuard organization. This table returns a single row with organization details.

**Required API Permission:** `Platform`

## Example queries

**Get organization details:**

```sql
select
  primary_hostname,
  name,
  automated_score,
  category_scores
from
  upguard_organisation;
```

**Get organization category score breakdown:**

```sql
select
  primary_hostname,
  name,
  category_scores->>'website_security' as website_security,
  category_scores->>'email_security' as email_security,
  category_scores->>'network_security' as network_security,
  category_scores->>'phishing_malware' as phishing_malware
from
  upguard_organisation;
```

**Get organization score and category scores:**

```sql
select
  primary_hostname,
  name,
  automated_score,
  category_scores
from
  upguard_organisation;
```
