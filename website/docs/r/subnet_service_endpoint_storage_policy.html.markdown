---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet_service_endpoint_storage_policy"
description: |-
  Manages a Subnet Service Endpoint Storage Policy.
---

# azurerm_subnet_service_endpoint_storage_policy

Manages a Subnet Service Endpoint Storage Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_subnet_service_endpoint_storage_policy" "example" {
  name                = "example-policy"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  definition {
    name        = "name1"
    description = "definition1"
    service     = "Microsoft.Storage"
    service_resources = [
      azurerm_resource_group.example.id,
      azurerm_storage_account.example.id
    ]
  }
  definition {
    name        = "name2"
    description = "definition2"
    service     = "Global"
    service_resources = [
      "/services/Azure",
      "/services/Azure/Batch",
      "/services/Azure/DataFactory",
      "/services/Azure/MachineLearning",
      "/services/Azure/ManagedInstance",
      "/services/Azure/WebPI",
    ]
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Subnet Service Endpoint Storage Policy. Changing this forces a new Subnet Service Endpoint Storage Policy to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Subnet Service Endpoint Storage Policy should exist. Changing this forces a new Subnet Service Endpoint Storage Policy to be created.

* `location` - (Required) The Azure Region where the Subnet Service Endpoint Storage Policy should exist. Changing this forces a new Subnet Service Endpoint Storage Policy to be created.

---

* `definition` - (Optional) A `definition` block as defined below

* `tags` - (Optional) A mapping of tags which should be assigned to the Subnet Service Endpoint Storage Policy.

---

A `definition` block supports the following:

* `name` - (Required) The name which should be used for this Subnet Service Endpoint Storage Policy Definition.

* `service` - (Optional) The type of service resources. Valid values are `Microsoft.Storage` or `Global`. When the `service_resources` property contains resource IDs, this property must be `Microsoft.Storage`. When the `service_resources` property contains Aliases, this property must be `Global`. Defaults to `Microsoft.Storage`. 

* `service_resources` - (Required) Specifies a list of resources or aliases that this Subnet Service Endpoint Storage Policy Definition applies to.

~> **Note:** The `service_resources` property must contain either Aliases or Resource IDs, but not both.

* `description` - (Optional) The description of this Subnet Service Endpoint Storage Policy Definition.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Subnet Service Endpoint Storage Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Subnet Service Endpoint Storage Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Subnet Service Endpoint Storage Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Subnet Service Endpoint Storage Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Subnet Service Endpoint Storage Policy.

## Import

Subnet Service Endpoint Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subnet_service_endpoint_storage_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/serviceEndpointPolicies/policy1
```
