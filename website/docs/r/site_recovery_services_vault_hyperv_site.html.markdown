---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_site_recovery_services_vault_hyperv_site"
description: |-
  Manages a HyperV Site in Recovery Service Vault.
---

# azurerm_site_recovery_services_vault_hyperv_site

Manages a HyperV Site in Recovery Service Vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "eastus"
}

resource "azurerm_recovery_services_vault" "example" {
  name                = "example-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  soft_delete_enabled = false
}


resource "azurerm_site_recovery_services_vault_hyperv_site" "example" {
  name              = "example-site"
  recovery_vault_id = azurerm_recovery_services_vault.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Recovery Service. Changing this forces a new Site to be created.

* `recovery_vault_id` - (Required) The ID of the Recovery Services Vault where the Site created. Changing this forces a new Site to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Recovery Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Recovery Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Service.
* `delete` - (Defaults to 3 hours) Used when deleting the Recovery Service.

## Import

Recovery Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_site_recovery_services_vault_hyperv_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1/replicationFabrics/fabric1
```
