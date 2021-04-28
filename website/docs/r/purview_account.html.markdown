---
subcategory: "Purview"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_purview_account"
description: |-
  Manages a Purview Account.
---

# azurerm_purview_account

Manages a Purview Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_purview_account" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Standard_4"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Purview Account should exist. Changing this forces a new Purview Account to be created.

* `name` - (Required) The name which should be used for this Purview Account. Changing this forces a new Purview Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Purview Account should exist. Changing this forces a new Purview Account to be created.

* `sku_name` - (Required) The SKU's capacity for platform size and catalog capabilities. Accepted values are `Standard_4` and `Standard_16`.

---

* `public_network_enabled` - (Optional) Should the Purview Account be visible to the public network? Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Purview Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Purview Account.

* `atlas_kafka_endpoint_primary_connection_string` - Atlas Kafka endpoint primary connection string.

* `atlas_kafka_endpoint_secondary_connection_string` - Atlas Kafka endpoint secondary connection string.

* `catalog_endpoint` - Catalog endpoint.

* `guardian_endpoint` - Guardian endpoint.

* `scan_endpoint` - Scan endpoint.

* `identity` - A `identity` block as defined below.

---

A `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The type of Managed Identity assigned to this Purview Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Purview Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Purview Account.
* `update` - (Defaults to 30 minutes) Used when updating the Purview Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Purview Account.

## Import

Purview Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_purview_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Purview/accounts/account1
```
