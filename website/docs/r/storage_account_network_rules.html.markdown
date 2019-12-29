---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_network_rules"
sidebar_current: "docs-azurerm-resource-storage-account-network-rules"
description: |-
  Manages network rules inside of a Azure Storage Account.
---

# azurerm_storage_account

Manages network rules inside of a Azure Storage Account.

~> **NOTE:** Network Rules can be defined either directly on the `azurerm_storage_account` resource, or using the `azurerm_storage_account_network_rules` resource - but the two cannot be used together. Spurious changes will occur if both are used against the same Storage Account. 

~> **NOTE:** Only one `azurerm_storage_account_network_rules` can be tied to an `azurerm_storage_account`. Spurious changes will occur if more than `azurerm_storage_account_network_rules` is tied to the same `azurerm_storage_account`.

~> **NOTE:** Deleting this resource updates the storage account back to the default values it had when the storage account was created.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = "${azurerm_resource_group.example.name}"
  location                 = "${azurerm_resource_group.example.location}"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = ["${azurerm_subnet.test.id}"]
  bypass = ["Metrics"]
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_name` - (Required) Specifies the name of the storage account. Changing this forces a new resource to be created. This must be unique across the entire Azure service, not just within the resource group.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage account. Changing this forces a new resource to be created.

* `default_action` - (Required) Specifies the default action of allow or deny when no other rules match. Valid options are `Deny` or `Allow`.

* `bypass` - (Optional)  Specifies whether traffic is bypassed for Logging/Metrics/AzureServices. Valid options are any combination of `Logging`, `Metrics`, `AzureServices`, or `None`. 

* `ip_rules` - (Optional) List of public IP or IP ranges in CIDR Format. Only IPV4 addresses are allowed. Private IP address ranges (as defined in [RFC 1918](https://tools.ietf.org/html/rfc1918#section-3)) are not allowed.

* `virtual_network_subnet_ids` - (Optional) A list of virtual network subnet ids to to secure the storage account.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account.

## Import

Storage Account Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_network_rules.storageAcc1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
