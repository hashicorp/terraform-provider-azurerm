---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_network_rules"
description: |-
  Manages network rules inside of a Azure Storage Account.
---

# azurerm_storage_account_network_rules

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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_network_rules" "test" {
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_name = azurerm_storage_account.test.name

  default_action             = "Allow"
  ip_rules                   = ["127.0.0.1"]
  virtual_network_subnet_ids = [azurerm_subnet.test.id]
  bypass                     = ["Metrics"]
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_name` - (Required) Specifies the name of the storage account. Changing this forces a new resource to be created. This must be unique across the entire Azure service, not just within the resource group.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage account. Changing this forces a new resource to be created.

* `default_action` - (Required) Specifies the default action of allow or deny when no other rules match. Valid options are `Deny` or `Allow`.

* `bypass` - (Optional)  Specifies whether traffic is bypassed for Logging/Metrics/AzureServices. Valid options are any combination of `Logging`, `Metrics`, `AzureServices`, or `None`.

-> **NOTE** User has to explicitly set `bypass` to empty slice (`[]`) to remove it.

* `ip_rules` - (Optional) List of public IP or IP ranges in CIDR Format. Only IPV4 addresses are allowed. Private IP address ranges (as defined in [RFC 1918](https://tools.ietf.org/html/rfc1918#section-3)) are not allowed.

-> **NOTE** Small address ranges using "/31" or "/32" prefix sizes are not supported. These ranges should be configured using individual IP address rules without prefix specified.

-> **NOTE** IP network rules have no effect on requests originating from the same Azure region as the storage account. Use Virtual network rules to allow same-region requests. Services deployed in the same region as the storage account use private Azure IP addresses for communication. Thus, you cannot restrict access to specific Azure services based on their public outbound IP address range.

-> **NOTE** User has to explicitly set `ip_rules` to empty slice (`[]`) to remove it.

* `virtual_network_subnet_ids` - (Optional) A list of virtual network subnet ids to to secure the storage account.

-> **NOTE** User has to explicitly set `virtual_network_subnet_ids` to empty slice (`[]`) to remove it.

* `private_link_access` - (Optional) One or More `private_link_access` block as defined below.

---

A `private_link_access` block supports the following:

* `endpoint_resource_id` - (Required) The resource id of the `azurerm_private_endpoint` of the resource access rule.

* `endpoint_tenant_id` - (Optional) The tenant id of the `azurerm_private_endpoint` of the resource access rule. Defaults to the current tenant id.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the  Network Rules for this Storage Account.
* `update` - (Defaults to 60 minutes) Used when updating the Network Rules for this Storage Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Rules for this Storage Account.
* `delete` - (Defaults to 60 minutes) Used when deleting the Network Rules for this Storage Account.

## Import

Storage Account Network Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_network_rules.storageAcc1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
