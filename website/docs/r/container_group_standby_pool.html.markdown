---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group_standby_pool"
description: |-
  Manages a Standby Pool for Container Group.
---

# azurerm_container_group_standby_pool

Manages a Standby Pool for Container Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_group" "example" {
  name                = "example-continst"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"

    ports {
      port     = 443
      protocol = "TCP"
    }
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }

  tags = {
    environment = "testing"
  }
}

resource "azurerm_container_group_standby_pool" "example" {
  name                = "example-pool"
  resource_group_name = azurerm_resource_group.example.name
  container_gorup_id  = azurerm_container_group.example.id
  max_ready_capacity  = 42
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Standby Pool for Container Group. Changing this forces a new Standby Pool for Container Group to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Standby Pool for Container Group should exist. Changing this forces a new Standby Pool for Container Group to be created.

* `container_gorup_id` - (Required) The ID of the Container Group of Standby Container Groups. Changing this forces a new Standby Pool for Container Group to be created.

* `max_ready_capacity` - (Required) The maximum number of Standby Container Groups in the Standby Pool.

---

* `container_group_revision` - (Optional) Specifies the Revision of Container Group.

* `refill_policy` - (Optional) Specifies the refill policy the standby pool is automatically refilled to maintain maxReadyCapacity. Possible value is `always`.

* `subnet_ids` - (Optional) Specifies a list of ID of subnets. Changing this forces a new Standby Pool for Container Group to be created.

* `zone` - (Optional) Specifies zones of Standby Container Group Pools.

* `tags` - (Optional) A mapping of tags which should be assigned to the Standby Pool for Container Group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Standby Pool for Container Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Standby Pool for Container Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Standby Pool for Container Group.
* `update` - (Defaults to 30 minutes) Used when updating the Standby Pool for Container Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Standby Pool for Container Group.

## Import

Standby Pool for Container Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_group_standby_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StandbyPool/standbyContainerGroupPools/pool1
```

