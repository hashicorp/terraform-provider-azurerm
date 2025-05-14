---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_eventhub_namespace"
description: |-
  Manages an EventHub Namespace.
---

# azurerm_eventhub_namespace

Manages an EventHub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example-namespace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 2

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the EventHub Namespace resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the namespace. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Defines which tier to use. Valid options are `Basic`, `Standard`, and `Premium`. Please note that setting this field to `Premium` will force the creation of a new resource.

* `capacity` - (Optional) Specifies the Capacity / Throughput Units for a `Standard` SKU namespace. Default capacity has a maximum of `2`, but can be increased in blocks of 2 on a committed purchase basis. Defaults to `1`.

* `auto_inflate_enabled` - (Optional) Is Auto Inflate enabled for the EventHub Namespace?

* `dedicated_cluster_id` - (Optional) Specifies the ID of the EventHub Dedicated Cluster where this Namespace should created. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `maximum_throughput_units` - (Optional) Specifies the maximum number of throughput units when Auto Inflate is Enabled. Valid values range from `1` - `20`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `network_rulesets` - (Optional) A `network_rulesets` block as defined below.

* `local_authentication_enabled` - (Optional) Is SAS authentication enabled for the EventHub Namespace? Defaults to `true`.

* `public_network_access_enabled` - (Optional) Is public network access enabled for the EventHub Namespace? Defaults to `true`.

* `minimum_tls_version` - (Optional) The minimum supported TLS version for this EventHub Namespace. Valid values are: `1.0`, `1.1` and `1.2`. Defaults to `1.2`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Event Hub Namespace. Possible values are `SystemAssigned` or `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this EventHub namespace.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

~> **Note:** Due to the limitation of the current Azure API, once an EventHub Namespace has been assigned an identity, it cannot be removed.

---

A `network_rulesets` block supports the following:

* `default_action` - (Required) The default action to take when a rule is not matched. Possible values are `Allow` and `Deny`.

* `public_network_access_enabled` - (Optional) Is public network access enabled for the EventHub Namespace? Defaults to `true`.

~> **Note:** The public network access setting at the network rule sets level should be the same as it's at the namespace level.

* `trusted_service_access_enabled` - (Optional) Whether Trusted Microsoft Services are allowed to bypass firewall.

* `virtual_network_rule` - (Optional) One or more `virtual_network_rule` blocks as defined below.

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

---

A `virtual_network_rule` block supports the following:

* `subnet_id` - (Required) The id of the subnet to match on.

* `ignore_missing_virtual_network_service_endpoint` - (Optional) Are missing virtual network service endpoints ignored? 

---

A `ip_rule` block supports the following:

* `ip_mask` - (Required) The IP mask to match on.

* `action` - (Optional) The action to take when the rule is matched. Possible values are `Allow`. Defaults to `Allow`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The EventHub Namespace ID.

* `identity` - An `identity` block as documented below.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization rule `RootManageSharedAccessKey`.

* `default_primary_connection_string_alias` - The alias of the primary connection string for the authorization rule `RootManageSharedAccessKey`, which is generated when disaster recovery is enabled.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string_alias` - The alias of the secondary connection string for the authorization rule `RootManageSharedAccessKey`, which is generated when disaster recovery is enabled.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace.

## Import

EventHub Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1
```
