---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_log_analytics_workspace_onboarding"
description: |-
  Manages a Security Insights Sentinel Onboarding States.
---

# azurerm_sentinel_log_analytics_workspace_onboarding

Manages a Security Insights Sentinel Onboarding.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "example" {
  resource_group_name          = azurerm_resource_group.example.name
  workspace_name               = azurerm_log_analytics_workspace.example.name
  customer_managed_key_enabled = false
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Security Insights Sentinel Onboarding States should exist. Changing this forces the Log Analytics Workspace off the board and onboard again.

* `workspace_name` - (Required) Specifies the Workspace Name. Changing this forces the Log Analytics Workspace off the board and onboard again. Changing this forces a new resource to be created.

* `customer_managed_key_enabled` - (Optional) Specifies if the Workspace is using Customer managed key. Defaults to `false`. Changing this forces a new resource to be created.

-> **Note:** To set up Microsoft Sentinel customer-managed key it needs to enable CMK on the workspace and add access policy to your Azure Key Vault. Details could be found on [this document](https://learn.microsoft.com/en-us/azure/sentinel/customer-managed-keys)

-> **Note:** Once a workspace is onboarded to Microsoft Sentinel with `customer_managed_key_enabled` set to true, it will not be able to be onboarded again with `customer_managed_key_enabled` set to false.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Security Insights Sentinel Onboarding States.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Insights Sentinel Onboarding States.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Insights Sentinel Onboarding States.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Insights Sentinel Onboarding States.

## Import

Security Insights Sentinel Onboarding States can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_log_analytics_workspace_onboarding.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/onboardingStates/defaults
```
