---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_network_rule"
description: |-
  Manages an Azure Storage Account Network Rule.
---

# azurerm_storage_account_network_rule

Manages an Azure Storage Account Network Rule.

~> **NOTE:** IP Network Rules, Virtual Network Rules and Resource Access Network Rules can be defined either directly on the [`azurerm_storage_account`](storage_account.html) resource, using the [`azurerm_storage_account_network_rules`](storage_account_network_rule.html) resource or using `azurerm_storage_account_network_rule` resources - but they cannot be used together. Spurious changes will occur if they are used against the same Storage Account.

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

  network_rules {
    default_action = "Deny"
    bypass         = ["AzureServices"]
  }

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_network_rule" "example" {
  storage_account_id = azurerm_storage_account.example.id
  ip_rule            = "127.0.0.1"
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Optional) Specifies the ID of the storage account. Changing this forces a new resource to be created.

* `ip_rule` - (Optional) Public IP or IP ranges in CIDR Format. Only IPV4 addresses are allowed. Private IP address ranges (as defined in [RFC 1918](https://tools.ietf.org/html/rfc1918#section-3)) are not allowed. Changing this forces a new resource to be created. This field cannot be specified if `virtual_network_rule` or `resource_access_rule` is specified.

-> **NOTE** Small address ranges using "/31" or "/32" prefix sizes are not supported. These ranges should be configured using individual IP address rules without prefix specified.

-> **NOTE** IP network rules have no effect on requests originating from the same Azure region as the storage account. Use Virtual network rules to allow same-region requests. Services deployed in the same region as the storage account use private Azure IP addresses for communication. Thus, you cannot restrict access to specific Azure services based on their public outbound IP address range.

* `virtual_network_rule` - (Optional) A virtual network subnet ids to secure the storage account. Changing this forces a new resource to be created. This field cannot be specified if `ip_rule` or `resource_access_rule` is specified.

* `resource_access_rule` - (Optional) A `resource_access_rule` block as defined below. Changing this forces a new resource to be created. This field cannot be specified if `ip_rule` or `virtual_network_rule` is specified.

---

A `resource_access_rule` block supports the following:

* `resource_id` - (Required) The resource id of the resource access rule to be granted access.

* `tenant_id` - (Optional) The tenant id of the resource of the resource access rule to be granted access. Defaults to the current tenant id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account Network Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Storage Account Network Rule.
* `update` - (Defaults to 10 minutes) Used when updating the Storage Account Network Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Network Rule.
* `delete` - (Defaults to 10 minutes) Used when deleting the Storage Account Network Rule.

## Import

Storage Account Network Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_network_rule.ip_rule /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount;ipAddressOrRange/127.0.0.1
terraform import azurerm_storage_account_network_rule.virtual_network_rule /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount;subnetId/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
terraform import azurerm_storage_account_network_rule.resource_access_rule /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount;tenantId/00000000-0000-0000-0000-000000000000/resourceId/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateEndpoints/myprivatelink
```
