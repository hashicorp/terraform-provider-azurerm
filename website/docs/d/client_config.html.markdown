---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_client_config"
description: |-
  Gets information about the configuration of the azurerm provider.
---

# Data Source: azurerm_client_config

Use this data source to access the configuration of the AzureRM provider.

## Example Usage

```hcl
data "azurerm_client_config" "current" {
}

output "account_id" {
  value = data.azurerm_client_config.current.client_id
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

* `client_id` is set to the Azure Client ID (Application Object ID).
* `tenant_id` is set to the Azure Tenant ID.
* `subscription_id` is set to the Azure Subscription ID.
* `object_id` is set to the Azure Object ID.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the client config.
