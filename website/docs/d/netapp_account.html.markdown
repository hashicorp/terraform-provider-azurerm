---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
description: |-
  Gets information about an existing NetApp Account
---

# Data Source: azurerm_netapp_account

Uses this data source to access information about an existing NetApp Account.

## NetApp Account Usage

```hcl
data "azurerm_netapp_account" "example" {
  resource_group_name = "acctestRG"
  name                = "acctestnetappaccount"
}

output "netapp_account_id" {
  value = data.azurerm_netapp_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of the NetApp Account.

* `resource_group_name` - The Name of the Resource Group where the NetApp Account exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Account exists.

* `nfsv4_id_domain` - The NFSv4.1 ID domain used for user and group name to UID and GID mappings.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Account.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.NetApp` - 2025-06-01
