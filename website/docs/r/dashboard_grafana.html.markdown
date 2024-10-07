---
subcategory: "Dashboard"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dashboard_grafana"
description: |-
  Manages a Dashboard Grafana.
---

# azurerm_dashboard_grafana

Manages a Dashboard Grafana.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dashboard_grafana" "example" {
  name                              = "example-dg"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = "West Europe"
  grafana_major_version             = 10
  api_key_enabled                   = true
  deterministic_outbound_ip_enabled = true
  public_network_access_enabled     = false

  identity {
    type = "SystemAssigned"
  }

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Dashboard Grafana. Changing this forces a new Dashboard Grafana to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Dashboard Grafana should exist. Changing this forces a new Dashboard Grafana to be created.

* `location` - (Required) Specifies the Azure Region where the Dashboard Grafana should exist. Changing this forces a new Dashboard Grafana to be created.

* `grafana_major_version` - (Required) Which major version of Grafana to deploy. Possible values are `9`, `10`. Changing this forces a new resource to be created.

* `api_key_enabled` - (Optional) Whether to enable the api key setting of the Grafana instance. Defaults to `false`.

* `auto_generated_domain_name_label_scope` - (Optional) Scope for dns deterministic name hash calculation. The only possible value is `TenantReuse`. Defaults to `TenantReuse`.

* `deterministic_outbound_ip_enabled` - (Optional) Whether to enable the Grafana instance to use deterministic outbound IPs. Defaults to `false`.

* `smtp` - (Optional) A `smtp` block as defined below.

* `azure_monitor_workspace_integrations` - (Optional) A `azure_monitor_workspace_integrations` block as defined below.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Dashboard Grafana to be created.

* `public_network_access_enabled` - (Optional) Whether to enable traffic over the public interface. Defaults to `true`.

* `sku` - (Optional) The name of the SKU used for the Grafana instance. Possible values are `Standard` and `Essential`. Defaults to `Standard`. Changing this forces a new Dashboard Grafana to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Dashboard Grafana.

* `zone_redundancy_enabled` - (Optional) Whether to enable the zone redundancy setting of the Grafana instance. Defaults to `false`. Changing this forces a new Dashboard Grafana to be created.

---

A `smtp` block supports the following:

* `enabled` - (Optional) Whether to enable the smtp setting of the Grafana instance. Defaults to `false`.

* `host` - (Required) SMTP server hostname with port, e.g. test.email.net:587

* `user` - (Required) User of SMTP authentication.

* `password` - (Required) Password of SMTP authentication.

* `start_tls_policy` - (Required) Whether to use TLS when connecting to SMTP server. Possible values are `OpportunisticStartTLS`, `NoStartTLS`, `MandatoryStartTLS`.

* `from_address` - (Required) Address used when sending emails.

* `from_name` - (Optional) Name used when sending emails. Defaults to `Azure Managed Grafana Notification`.

* `verification_skip_enabled` - (Optional) Whether verify SSL for SMTP server. Defaults to `false`.

---

An `azure_monitor_workspace_integrations` block supports the following:

* `resource_id` - (Required) Specifies the resource ID of the connected Azure Monitor Workspace.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Optional) Specifies the list of User Assigned Managed Service Identity IDs which should be assigned to this Dashboard Grafana. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dashboard Grafana.

* `endpoint` - The endpoint of the Grafana instance.

* `grafana_version` - The full Grafana software semantic version deployed.

* `identity` - An `identity` block as defined below.

* `outbound_ip` - List of outbound IPs if deterministicOutboundIP is enabled.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dashboard Grafana.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dashboard Grafana.
* `update` - (Defaults to 30 minutes) Used when updating the Dashboard Grafana.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dashboard Grafana.

## Import

Dashboard Grafana can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dashboard_grafana.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Dashboard/grafana/workspace1
```
