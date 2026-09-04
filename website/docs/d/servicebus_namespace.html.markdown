---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_servicebus_namespace"
description: |-
  Gets information about an existing ServiceBus Namespace.
---

# Data Source: azurerm_servicebus_namespace

Use this data source to access information about an existing ServiceBus Namespace.

## Example Usage

```hcl
data "azurerm_servicebus_namespace" "example" {
  name                = "examplenamespace"
  resource_group_name = "example-resources"
}

output "location" {
  value = data.azurerm_servicebus_namespace.example.location
}
```

## Arguments Reference

* `name` - Specifies the name of the ServiceBus Namespace.

* `resource_group_name` - Specifies the name of the Resource Group where the ServiceBus Namespace exists.

## Attributes Reference

* `location` - The location of the Resource Group in which the ServiceBus Namespace exists.

* `identity` - An `identity` block as defined below.

* `sku` - The Tier used for the ServiceBus Namespace.

* `capacity` - The capacity of the ServiceBus Namespace.

* `premium_messaging_partitions` - The messaging partitions of the ServiceBus Namespace.

* `customer_managed_key` - A `customer_managed_key` block as defined below.

* `local_auth_enabled` - Whether or not SAS authentication is enabled for the Service Bus Namespace.

* `public_network_access_enabled` - Whether public network access is enabled for the Service Bus Namespace.

* `minimum_tls_version` - The minimum supported TLS version for this Service Bus Namespace.

* `endpoint` - The URL to access the ServiceBus Namespace.

* `network_rule_set` - A `network_rule_set` block as defined below.

* `tags` - A mapping of tags assigned to the resource.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity configured on this Service Bus Namespace.

* `identity_ids` - The User Assigned Managed Identity IDs assigned to this Service Bus Namespace.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Service Bus Namespace.

---

A `customer_managed_key` block exports the following:

* `key_vault_key_id` - The ID of the Key Vault Key used to encrypt data in this Service Bus Namespace.

* `identity_id` - The ID of the User Assigned Identity that has access to the key.

* `infrastructure_encryption_enabled` - Whether Infrastructure Encryption (Double Encryption) is enabled.

---

A `network_rule_set` block exports the following:

* `default_action` - The default action for the Network Rule Set.

* `public_network_access_enabled` - Whether traffic over public network is allowed by the Network Rule Set.

* `trusted_services_allowed` - Whether trusted Azure services are allowed to bypass firewall configuration.

* `ip_rules` - The IP Addresses or CIDR Blocks which are able to access the Service Bus Namespace.

* `network_rules` - One or more `network_rules` blocks as defined below.

---

A `network_rules` block exports the following:

* `subnet_id` - The Subnet ID which is able to access this Service Bus Namespace.

* `ignore_missing_vnet_service_endpoint` - Whether the Network Rule Set ignores missing Virtual Network Service Endpoint configuration in the Subnet.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the ServiceBus Namespace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ServiceBus` - 2024-01-01
