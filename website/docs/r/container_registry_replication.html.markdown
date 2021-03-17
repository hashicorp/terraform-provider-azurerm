---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_replication"
description: |-
  Manages a Container Registry Replication.
---

# azurerm_container_registry_replication

Manages a Container Registry Replication.

~> **NOTE on Container Registry and Container Registry Replication's:** Terraform currently
provides both a standalone [Container Registry Replication](container_registry_replication.html), and allows for replications to be defined in-line within the [Container Registry](container_registry.html).
At this time you cannot use a Container Registry with in-line Replication in conjunction with any Container Registry Replication resources. Doing so will cause a conflict of Replication configurations.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "example" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
}

resource "azurerm_container_registry_replication" "example" {
  name                = "myreplication"
  resource_group_name = azurerm_resource_group.example.name
  registry_name       = azurerm_container_registry.example.name
  location            = "West US"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Container Registry Replication. Changing this forces a new Container Registry Replication to be created.

* `registry_name` - (Required) The Name of Container registry this Webhook belongs to. Changing this forces a new Container Registry Replication to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Container Registry Replication should exist. Changing this forces a new Container Registry Replication to be created.

* `location` - (Required) The Azure Region where the Container Registry Replication should exist. Changing this forces a new Container Registry Replication to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Container Registry Replication.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Container Registry Replication.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Replication.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Replication.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Replication.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Replication.

## Import

Container Registry Replications can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_replication.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1/replications/myreplication1
```
