---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advanced_threat_protection"
description: |-
  Manages a resources Advanced Threat Protection setting.
---

# azurerm_advanced_threat_protection

Manages a resources Advanced Threat Protection setting.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "atp-example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                = "examplestorage"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "example"
  }
}

resource "azurerm_advanced_threat_protection" "example" {
  target_resource_id = azurerm_storage_account.example.id
  enabled            = true
}
```

## Argument Reference

The following arguments are supported:

* `target_resource_id` - (Required) The ID of the Azure Resource which to enable Advanced Threat Protection on. Changing this forces a new resource to be created.

* `enabled` - (Required) Should Advanced Threat Protection be enabled on this resource?


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Advanced Threat Protection resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Advanced Threat Protection.
* `update` - (Defaults to 30 minutes) Used when updating the Advanced Threat Protection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Advanced Threat Protection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Advanced Threat Protection.

## Import

Advanced Threat Protection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_advanced_threat_protection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleResourceGroup/providers/Microsoft.Storage/storageAccounts/exampleaccount/providers/Microsoft.Security/advancedThreatProtectionSettings/default
```
