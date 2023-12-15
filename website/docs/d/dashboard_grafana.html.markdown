---
subcategory: "Dashboard"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dashboard_grafana"
description: |-
  Gets information about an existing Grafana Dashboard.
---

# Data Source: azurerm_dashboard_grafana

Use this data source to access information about an existing Grafana Dashboard.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_dashboard_grafana" "example" {
  name                = "example-grafana-dashboard"
  resource_group_name = "example-rg"
}

output "name" {
  value = data.azurerm_dashboard_grafana.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the grafana dashboard.

* `resource_group_name` - (Required) Name of the resource group where resource belongs to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Grafana Dashboard.

* `location` - Azure location where the resource exists.

* `api_key_enabled` - Whether the api key setting of the Grafana instance is enabled.

* `auto_generated_domain_name_label_scope` - Scope for dns deterministic name hash calculation.

* `deterministic_outbound_ip_enabled` - Whether the Grafana instance uses deterministic outbound IPs.

* `grafana_major_version` - Major version of Grafana instance.

* `azure_monitor_workspace_integrations` - Integrations for Azure Monitor Workspace.

* `identity` - The managed identity of the grafana resource.

* `public_network_access_enabled` - Whether or not public endpoint access is allowed for this server.

* `endpoint` - The endpoint of the Grafana instance.

* `grafana_version` - The full Grafana software semantic version deployed.

* `outbound_ip` - List of outbound IPs if deterministicOutboundIP is enabled.

* `sku` - The name of the SKU used for the Grafana instance.

* `zone_redundancy_enabled` - The zone redundancy setting of the Grafana instance.

* `tags` - A mapping of tags to assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API.
