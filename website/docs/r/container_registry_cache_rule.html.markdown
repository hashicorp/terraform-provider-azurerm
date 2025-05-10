---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_cache_rule"
description: |-
  Manages an Azure Container Registry Cache Rule.

---

# azurerm_container_registry_cache_rule

Manages an Azure Container Registry Cache Rule.

~> **Note:** All arguments including the access key will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_registry" "acr" {
  name                = "containerRegistry1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}

resource "azurerm_container_registry_cache_rule" "cache_rule" {
  name                  = "cacherule"
  container_registry_id = azurerm_container_registry.acr.id
  target_repo           = "target"
  source_repo           = "docker.io/hello-world"
  credential_set_id     = "${azurerm_container_registry.acr.id}/credentialSets/example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Registry Cache Rule. Only Alphanumeric characters allowed. Changing this forces a new resource to be created.

* `container_registry_id` - (Required) The ID of the Container Registry where the Cache Rule should apply. Changing this forces a new resource to be created.

* `source_repo` - (Required) The name of the source repository path. Changing this forces a new resource to be created. 

* `target_repo` - (Required) The name of the new repository path to store artifacts. Changing this forces a new resource to be created.

* `credential_set_id` - (Optional) The ARM resource ID of the Credential Store which is associated with the Cache Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry Cache Rule.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Cache Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Cache Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Cache Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Cache Rule.

## Import

Container Registry Cache Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_cache_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.ContainerRegistry/registries/myRegistry/cacheRules/myCacheRule
```
