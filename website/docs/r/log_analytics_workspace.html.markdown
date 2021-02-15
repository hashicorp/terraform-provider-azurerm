---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_workspace"
description: |-
  Manages a Log Analytics (formally Operational Insights) Workspace.
---

# azurerm_log_analytics_workspace

Manages a Log Analytics (formally Operational Insights) Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Log Analytics Workspace. Workspace name should include 4-63 letters, digits or '-'. The '-' shouldn't be the first or the last symbol. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics workspace is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Optional) Specifies the Sku of the Log Analytics Workspace. Possible values are `Free`, `PerNode`, `Premium`, `Standard`, `Standalone`, `Unlimited`, and `PerGB2018` (new Sku as of `2018-04-03`). Defaults to `PerGB2018`.

~> **NOTE:** A new pricing model took effect on `2018-04-03`, which requires the SKU `PerGB2018`. If you're provisioned resources before this date you have the option of remaining with the previous Pricing SKU and using the other SKU's defined above. More information about [the Pricing SKU's is available at the following URI](http://aka.ms/PricingTierWarning).

~> **NOTE:** The `Free` SKU has a default `daily_quota_gb` value of `0.5` (GB).

* `retention_in_days` - (Optional) The workspace data retention in days. Possible values are either 7 (Free Tier only) or range between 30 and 730.

* `daily_quota_gb` - (Optional) The workspace daily quota for ingestion in GB.  Defaults to -1 (unlimited) if omitted.

~> **NOTE:** When `sku_name` is set to `Free` this field can be set to a maximum of `0.5` (GB), and has a default value of `0.5`. 

* `internet_ingestion_enabled ` - (Optional) Should the Log Analytics Workflow support ingestion over the Public Internet? Defaults to `true`.

* `internet_query_enabled` - (Optional) Should the Log Analytics Workflow support querying over the Public Internet? Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **NOTE:** If a `azurerm_log_analytics_workspace` is connected to a `azurerm_log_analytics_cluster` via a `azurerm_log_analytics_linked_service` it will not be able to be modified until link between the workspace and the cluster has been broken by deleting the `azurerm_log_analytics_linked_service` resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Workspace ID.

* `primary_shared_key` - The Primary shared key for the Log Analytics Workspace.

* `secondary_shared_key` - The Secondary shared key for the Log Analytics Workspace.

* `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.

## Import

Log Analytics Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1
```
