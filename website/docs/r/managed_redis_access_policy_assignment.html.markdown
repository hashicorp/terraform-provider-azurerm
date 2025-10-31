---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_access_policy_assignment"
description: |-
  Manages a Managed Redis Access Policy Assignment.
---

# azurerm_managed_redis_access_policy_assignment

Manages a Managed Redis Access Policy Assignment.

## Example Usage

```hcl
data "azuread_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_redis" "example" {
  name                = "example-managedredis"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Enterprise_E10"

  default_database {
    access_keys_authentication_enabled = true
  }
}

resource "azurerm_managed_redis_access_policy_assignment" "example" {
  name               = "example
  managed_redis_id   = azurerm_managed_redis.example.id
  access_policy_name = "Data Contributor"
  object_id          = data.azuread_client_config.current.object_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Access Policy Assignment. Changing this forces a new Access Policy Assignment to be created.

* `managed_redis_id` - (Required) The ID of the Managed Redis instance. Changing this forces a new Access Policy Assignment to be created.

* `access_policy_name` - (Required) The name of the Access Policy to be assigned. Changing this forces a new Redis Cache Access Policy Assignment to be created.

* `object_id` - (Required) The object ID of the Azure Active Directory user, group, service principal, or managed identity to assign the access policy to. Changing this forces a new Access Policy Assignment to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Access Policy Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Managed Redis Access Policy Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Access Policy Assignment.
* `delete` - (Defaults to 5 minutes) Used when deleting the Managed Redis Access Policy Assignment.

## Import

Managed Redis Access Policy Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_redis_access_policy_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Cache/redisEnterprise/redis1/databases/default/accessPolicyAssignments/assignment1
```
