---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_network_rule_set"
description: |-
  Manages an Azure Container Registry Network rule set.

---

# azurerm_container_registry_network_rule_set

Manages an Azure Container Registry Network rule set.

~> **NOTE:** It's possible to define Container Registry Network rule set both within [the `azurerm_container_registry` resource](container_registry.html) via the `network_rule_set` block and by using [the `azurerm_container_registry_network_rule_set` resource](container_registry_network_rule_set.html). However it's not possible to use both methods to manage Network rule set within a Container Registry, since there'll be conflicts.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "resourcegroup1"
  location = "East US"
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  sku                 = "Premium"
}

resource "azurerm_virtual_network" "vnet" {
  name                = "virtual-network"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "subnet" {
  name                 = "acr-subnet"
  virtual_network_name = azurerm_virtual_network.vnet.name
  resource_group_name  = azurerm_resource_group.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
  service_endpoints    = ["Microsoft.ContainerRegistry"]
}

resource "azurerm_container_registry_network_rule_set" "ruleset" {
  resource_group_name     = azurerm_container_registry.acr.resource_group_name
  container_registry_name = azurerm_container_registry.acr.name
  network_rule_set {
    default_action = "Deny"
    ip_rule {
      action   = "Allow"
      ip_range = "43.0.0.0/24"
    }
    virtual_network {
      action    = "Allow"
      subnet_id = azurerm_subnet.subnet.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `container_registry_name` - (Required) Specifies the name of the Container Registry.

* `resource_group_name` - (Required) The name of the resource group in which the Container Registry is existing.

* `network_rule_set` - A `network_rule_set` block as documented below.

`network_rule_set` supports the following:

* `default_action` - (Optional) The behaviour for requests matching no rules. Either `Allow` or `Deny`. Defaults to `Allow`

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

* `virtual_network` - (Optional) One or more `virtual_network` blocks as defined below.

~> **NOTE:** `network_rule_set ` is only supported with the `Premium` SKU at this time.

~> **NOTE:** Azure automatically configures Network Rules - to remove these you'll need to specify an `network_rule_set` block with `default_action` set to `Deny`.

`ip_rule` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `ip_range` - (Required) The CIDR block from which requests will match the rule.

`virtual_network` supports the following:

* `action` - (Required) The behaviour for requests matching this rule. At this time the only supported value is `Allow`

* `subnet_id` - (Required) The subnet id from which requests will match the rule.


---
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Container Registry Network rule set.

-> **NOTE:** This Identifier is unique to Terraform and doesn't map to an existing object within Azure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Container Registry Network rule set.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Container Registry Network rule set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Container Registry Network rule set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Container Registry Network rule set.

## Import

Azure Container Registry Network rule set can be imported using the Resource ID of the Container Registry, e.g.

```shell
terraform import azurerm_container_registry_network_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1
```
