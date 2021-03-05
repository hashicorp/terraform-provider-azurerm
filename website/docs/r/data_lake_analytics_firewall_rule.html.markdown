---
subcategory: "Data Lake"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_lake_analytics_firewall_rule"
description: |-
  Manages a Azure Data Lake Analytics Firewall Rule.
---

# azurerm_data_lake_analytics_firewall_rule

Manages a Azure Data Lake Analytics Firewall Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex_datalake_fw_rule"
  location = "West Europe"
}

resource "azurerm_data_lake_store" "example" {
  name                = "tfexdatalakestore"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_data_lake_analytics_account" "example" {
  name                = "tfexdatalakeaccount"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  default_store_account_name = azurerm_data_lake_store.example.name
}

resource "azurerm_data_lake_analytics_firewall_rule" "example" {
  name                = "office-ip-range"
  account_name        = azurerm_data_lake_analytics.example.name
  resource_group_name = azurerm_resource_group.example.name
  start_ip_address    = "1.2.3.4"
  end_ip_address      = "2.3.4.5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Lake Analytics. Changing this forces a new resource to be created. Has to be between 3 to 24 characters.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Lake Analytics.

* `account_name` - (Required) Specifies the name of the Data Lake Analytics for which the Firewall Rule should take effect.

* `start_ip_address` - (Required) The Start IP address for the firewall rule.

* `end_ip_address` - (Required) The End IP Address for the firewall rule.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Data Lake Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Lake Analytics Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Data Lake Analytics Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Lake Analytics Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Lake Analytics Firewall Rule.

## Import

Data Lake Analytics Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_lake_analytics_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DataLakeAnalytics/accounts/mydatalakeaccount/firewallRules/rule1
```
