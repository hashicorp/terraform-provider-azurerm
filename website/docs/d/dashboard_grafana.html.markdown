---
subcategory: "Dashboard"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dashboard_grafana"
description: |-
  Get information of a Dashboard Grafana.
---

# Data Source: azurerm_dashboard_grafana

Use this data to access information about an existing Dashboard Grafana.

## Example Usage

```hcl
data "azurerm_dashboard_grafana" "example" {
  name                = "example-dg"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the Grafana dashboard.

- `resource_group_name` - (Required) The name of the resource group in which the Grafana dashboard is located.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `api_key_enabled` - Whether the api key setting of the Grafana instance is enabled.

- `auto_generated_domain_name_label_scope` - Scope for dns deterministic name hash calculation. 

- `deterministic_outbound_ip_enabled` - Whether the Grafana instance is enabled to use deterministic outbound IPs.

- `azure_monitor_workspace_integrations` - A `azure_monitor_workspace_integrations` block as defined below.

- `location` - The Azure Region where the Grafana dashboard is located.

- `public_network_access_enabled` - Whether the Grafana instance is enabled to use public network access.

- `sku` - The SKU of the Grafana instance.

- `zone_redundancy_enabled` - Whether the Grafana instance is enabled to use zone redundancy.

- `endpoint` - The endpoint of the Grafana instance.

- `grafana_major_version` - The major version of Grafana deployed.

- `grafana_version` - The full Grafana software semantic version deployed.

- `outbound_ip` - List of outbound IPs if deterministicOutboundIP is enabled.

- `tags` - A mapping of tags assigned to the resource.

---

A `azure_monitor_workspace_integrations` block exports the following:

- `resource_id` - The resource ID of the connected Azure Monitor Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

- `read` - (Defaults to 5 minutes) Used when retrieving the Dashboard Grafana.
