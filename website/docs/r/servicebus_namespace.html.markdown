---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
description: |-
  Manages a ServiceBus Namespace.
---

# azurerm_servicebus_namespace

Manages a ServiceBus Namespace.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "terraform-servicebus"
  location = "West Europe"
}

resource "azurerm_servicebus_namespace" "example" {
  name                = "tfex-servicebus-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  tags = {
    source = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Service Bus Namespace resource . Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to Changing this forces a new resource to be created.
    create the namespace.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which tier to use. Options are `Basic`, `Standard` or `Premium`. Please note that setting this field to `Premium` will force the creation of a new resource.

* `identity` - (Optional) An `identity` block as defined below.

* `capacity` - (Optional) Specifies the capacity. When `sku` is `Premium`, capacity can be `1`, `2`, `4`, `8` or `16`. When `sku` is `Basic` or `Standard`, capacity can be `0` only.

* `premium_messaging_partitions` - (Optional) Specifies the number messaging partitions. Only valid when `sku` is `Premium` and the minimum number is `1`. Possible values include `0`, `1`, `2`, and `4`. Defaults to `0` for Standard, Basic namespace. Changing this forces a new resource to be created.

-> **Note:** It's not possible to change the partitioning option on any existing namespace. The number of partitions can only be set during namespace creation. Please check the doc https://learn.microsoft.com/en-us/azure/service-bus-messaging/enable-partitions-premium for more feature restrictions.

* `customer_managed_key` - (Optional) An `customer_managed_key` block as defined below.

* `local_auth_enabled` - (Optional) Whether or not SAS authentication is enabled for the Service Bus namespace. Defaults to `true`.

* `public_network_access_enabled` - (Optional) Is public network access enabled for the Service Bus Namespace? Defaults to `true`.

* `minimum_tls_version` - (Optional) The minimum supported TLS version for this Service Bus Namespace. Valid values are: `1.0`, `1.1` and `1.2`. Defaults to `1.2`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `network_rule_set` - (Optional) An `network_rule_set` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Service Bus Namespace. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Service Bus namespace.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

-> **Note:** Once customer-managed key encryption has been enabled, it cannot be disabled.

-> **Note:** The `customer_managed_key` block should only be used for Service Bus Namespaces with a User Assigned identity. To create a Customer Managed Key for a Service Bus Namespace with a System Assigned identity, use the `azurerm_servicebus_namespace_customer_managed_key` resource and add `customer_managed_key` to `ignore_changes`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Service Bus Namespace.

* `identity_id` - (Required) The ID of the User Assigned Identity that has access to the key.

* `infrastructure_encryption_enabled` - (Optional) Used to specify whether enable Infrastructure Encryption (Double Encryption). Changing this forces a new resource to be created.

---

A `network_rule_set` block supports the following:

* `default_action` - (Optional) Specifies the default action for the Network Rule Set. Possible values are `Allow` and `Deny`. Defaults to `Allow`.

* `public_network_access_enabled` - (Optional) Whether to allow traffic over public network. Possible values are `true` and `false`. Defaults to `true`.

-> **Note:** To disable public network access, you must also configure the property `public_network_access_enabled`.

* `trusted_services_allowed` - (Optional) Are Azure Services that are known and trusted for this resource type are allowed to bypass firewall configuration? See [Trusted Microsoft Services](https://github.com/MicrosoftDocs/azure-docs/blob/master/articles/service-bus-messaging/includes/service-bus-trusted-services.md)

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Service Bus Namespace.

* `network_rules` - (Optional) One or more `network_rules` blocks as defined below.

---

A `network_rules` block supports the following:

* `subnet_id` - (Required) The Subnet ID which should be able to access this Service Bus Namespace.

* `ignore_missing_vnet_service_endpoint` - (Optional) Should the Service Bus Namespace Network Rule Set ignore missing Virtual Network Service Endpoint option in the Subnet? Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Service Bus Namespace ID.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Service Bus Namespace.

* `endpoint` - The URL to access the Service Bus Namespace.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

---

The following attributes are exported only if there is an authorization rule named `RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Bus Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Bus Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the Service Bus Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Bus Namespace.

## Import

Service Bus Namespace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_servicebus_namespace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ServiceBus/namespaces/sbns1
```
