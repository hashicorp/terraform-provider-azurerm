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

* `sku` - (Required) Defines which tier to use. Valid options are `Basic` and `Standard`.

* `capacity` - (Optional) Specifies the Capacity / Throughput Units for a `Standard` SKU namespace. Valid values range from `1` - `20`.

* `auto_inflate_enabled` - (Optional) Is Auto Inflate enabled for the EventHub Namespace?

* `maximum_throughput_units` - (Optional) Specifies the maximum number of throughput units when Auto Inflate is Enabled. Valid values range from `1` - `20`.

* `kafka_enabled` - (Optional / **Deprecated**) Is Kafka enabled for the EventHub Namespace? Defaults to `false`.

-> **NOTE:** `kafka_enabled` is now configured depending on the `sku` being provisioned, where this is Disabled for a `Basic` sku and Enabled for a Standard sku.  

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `network_rulesets` - (Optional) A `network_rulesets` block as defined below.

---

A `network_rulesets` block supports the following:

* `default_action` - (Required) The default action to take when a rule is not matched. Possible values are `Allow` and `Deny`.

* `virtual_network_rule` - (Optional) One or more `virtual_network_rule` blocks as defined below.

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

---

A `virtual_network_rule` block supports the following:

* `subnet_id` - (Required) The id of the subnet to match on.

* `ignore_missing_virtual_network_service_endpoint` - (Optional) Are missing virtual network service endpoints ignored? Defaults to `false`.

---

A `ip_rule` block supports the following:

* `ip_mask` - (Required) The ip mask to match on.

* `action` - (Optional) The action to take when the rule is  matched. Possible values are `Allow`.

## Attributes Reference

The following attributes are exported:

* `id` - The EventHub Namespace ID.

The following attributes are exported only if there is an authorization rule named
`RootManageSharedAccessKey` which is created automatically by Azure.

* `default_primary_connection_string` - The primary connection string for the authorization
    rule `RootManageSharedAccessKey`.

* `default_secondary_connection_string` - The secondary connection string for the
    authorization rule `RootManageSharedAccessKey`.

* `default_primary_key` - The primary access key for the authorization rule `RootManageSharedAccessKey`.

* `default_secondary_key` - The secondary access key for the authorization rule `RootManageSharedAccessKey`.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the EventHub Namespace.
* `update` - (Defaults to 30 minutes) Used when updating the EventHub Namespace.
* `read` - (Defaults to 5 minutes) Used when retrieving the EventHub Namespace.
* `delete` - (Defaults to 30 minutes) Used when deleting the EventHub Namespace.

## Import

EventHub Namespaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_eventhub_namespace.namespace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.EventHub/namespaces/namespace1
```
