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
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
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

---

* `allow_resource_only_permissions` - (Optional) Specifies if the log Analytics Workspace allow users accessing to data associated with resources they have permission to view, without permission to workspace. Defaults to `true`.

* `local_authentication_disabled` - (Optional) Specifies if the log Analytics workspace should enforce authentication using Azure AD. Defaults to `false`.

* `sku` - (Optional) Specifies the SKU of the Log Analytics Workspace. Possible values are `PerNode`, `Premium`, `Standard`, `Standalone`, `Unlimited`, `CapacityReservation`, `PerGB2018`, and `LACluster`. Defaults to `PerGB2018`.

~> **Note:** `sku` should only be set to `LACluster` when the Log Analytics Workspace is linked to a Log Analytics Cluster. Additionally, `sku` cannot be modified while linked.

~> **Note:** Changing `sku` forces a new Log Analytics Workspace to be created, except when changing between `PerGB2018` and `CapacityReservation`. Changing `sku` to `CapacityReservation` or changing `reservation_capacity_in_gb_per_day` to a higher tier will lead to a 31-days commitment period, during which the SKU cannot be changed to a lower one. Please refer to [official documentation](https://learn.microsoft.com/en-us/azure/azure-monitor/logs/cost-logs#commitment-tiers) for further information.

-> **Note:** A new pricing model took effect on `2018-04-03`, which requires the SKU `PerGB2018`. If you've provisioned resources before this date you have the option of remaining with the previous Pricing SKU and using the other SKUs defined above. More information about [the Pricing SKUs is available at the following URI](https://aka.ms/PricingTierWarning).

* `retention_in_days` - (Optional) The workspace data retention in days. Possible values are between `30` and `730`.

* `daily_quota_gb` - (Optional) The workspace daily quota for ingestion in GB. Defaults to `-1` (unlimited) if omitted.

* `cmk_for_query_forced` - (Optional) Is Customer Managed Storage mandatory for query management?

* `identity` - (Optional) An `identity` block as defined below.

* `internet_ingestion_enabled` - (Optional) Should the Log Analytics Workspace support ingestion over the Public Internet? Defaults to `true`.

* `internet_query_enabled` - (Optional) Should the Log Analytics Workspace support querying over the Public Internet? Defaults to `true`.

* `reservation_capacity_in_gb_per_day` - (Optional) The capacity reservation level in GB for this workspace. Possible values are `100`, `200`, `300`, `400`, `500`, `1000`, `2000` and `5000`.

~> **Note:** `reservation_capacity_in_gb_per_day` can only be used when the `sku` is set to `CapacityReservation`.

* `data_collection_rule_id` - (Optional) The ID of the Data Collection Rule to use for this workspace.

* `immediate_data_purge_on_30_days_enabled` - (Optional) Whether to remove the data in the Log Analytics Workspace immediately after 30 days.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Log Analytics Workspace. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you) and `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field.

-> **Note:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the Log Analytics Workspace has been created.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Log Analytics Workspace ID.

* `primary_shared_key` - The Primary shared key for the Log Analytics Workspace.

* `secondary_shared_key` - The Secondary shared key for the Log Analytics Workspace.

* `workspace_id` - The Workspace (or Customer) ID for the Log Analytics Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.

## Import

Log Analytics Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_workspace.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1
```
