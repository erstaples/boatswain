# Global Values

Global values are useful for defining ecosystem-wide variables. For instance, let's say you have a New Relic account and you'd like to monitor multiple applications with it. Rather than repeating the license key over and over across multiple Helm packages, you can define it once in the global values.

```yaml
Global:
  NewRelic:
    LicenseKey: abc123
```

Global values operate using the same logic as standard Helm package values files. `.globals/values.prod.yaml` values will override the `.global/values.yaml` values for production, `values.staging.yaml` for staging, `values.env.yaml` for local deployments.

Global values can be accessed via the `.Values.Global` value object. For example, in the Hydra template, we use the New Relic global property like so:

```yaml
- name: NEWRELIC_LICENSE_KEY
  value: "{{ .Values.Global.NewRelic.LicenseKey }}"
```
