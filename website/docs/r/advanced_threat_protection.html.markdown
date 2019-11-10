---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advanced_threat_protection"
sidebar_current: "docs-azurerm-resource-azurerm-advanced-threat-protection-x"
description: |-
  Manages a resources Advanced Threat Protection setting.
---

# azurerm_advanced_threat_protection

Manages a resources Advanced Threat Protection setting..

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "atp-example"
  location = "northeurope"
}

resource "azurerm_storage_account" "example" {
  name                = "examplestorage"
  resource_group_name = "${azurerm_resource_group.example.name}"

  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "example"
  }
}

resource "azurerm_advanced_threat_protection" "example" {
  target_resource_id = "${azurerm_storage_account.example.id}"
  enabled            = true
}

```

## Argument Reference

The following arguments are supported:

* `target_resource_id` - (Required) The azure resource ID of the resource to enable the setting on. Changing this forces a new resource to be created.

* `enabled` - (Required) Whether to enable or disable Advanced Threat Protection on this resource.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Advanced Threat Protection resource.


## Import

Analysis Services Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_advanced_threat_protection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleResourceGroup/providers/Microsoft.Storage/storageAccounts/exampleaccount/providers/Microsoft.Security/advancedThreatProtectionSettings/default
```
