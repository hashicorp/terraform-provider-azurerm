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

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Purview Account should exist. Changing this forces a new Purview Account to be created.

* `identity` - (Required) An `identity` block as defined below.

* `name` - (Required) The name which should be used for this Purview Account. Changing this forces a new Purview Account to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Purview Account should exist. Changing this forces a new Purview Account to be created.

---

* `public_network_enabled` - (Optional) Should the Purview Account be visible to the public network? Defaults to `true`.

* `managed_resource_group_name` - (Optional) The name which should be used for the new Resource Group where Purview Account creates the managed resources. Changing this forces a new Purview Account to be created.

~> **Note:** `managed_resource_group_name` must be a new Resource Group.

* `tags` - (Optional) A mapping of tags which should be assigned to the Purview Account.

---

The `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Purview Account. Possible values are `UserAssigned` and `SystemAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Purview Account.

~> **Note:** This is required when `type` is set to `UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Purview Account.

* `atlas_kafka_endpoint_primary_connection_string` - Atlas Kafka endpoint primary connection string.

* `atlas_kafka_endpoint_secondary_connection_string` - Atlas Kafka endpoint secondary connection string.

* `catalog_endpoint` - Catalog endpoint.

* `guardian_endpoint` - Guardian endpoint.

* `scan_endpoint` - Scan endpoint.

* `identity` - A `identity` block as defined below.

* `managed_resources` - A `managed_resources` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `managed_resources` block exports the following:

* `event_hub_namespace_id` - The ID of the managed event hub namespace.

* `resource_group_id` - The ID of the managed resource group.

* `storage_account_id` - The ID of the managed storage account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Purview Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Purview Account.
* `update` - (Defaults to 30 minutes) Used when updating the Purview Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Purview Account.

## Import

Purview Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_purview_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Purview/accounts/account1
```
