---
subcategory: "Managed Redis"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_redis_access_policy_assignment"
description: |-
  Gets information about an existing Managed Redis Access Policy Assignment.
---

# Data Source: azurerm_managed_redis_access_policy_assignment

Use this data source to access information about an existing Managed Redis Access Policy Assignment.

## Example Usage

```hcl
data "azurerm_managed_redis_access_policy_assignment" "example" {
  name                = "example-assignment"
  managed_redis_name  = "example-managedredis"
  resource_group_name = "example-resources"
}

output "access_policy_name" {
  value = data.azurerm_managed_redis_access_policy_assignment.example.access_policy_name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Access Policy Assignment.

* `managed_redis_name` - (Required) The name of the Managed Redis instance.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Redis instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Access Policy Assignment.

* `access_policy_name` - The name of the access policy assigned.

* `object_id` - The object ID of the Azure Active Directory user, group, service principal, or managed identity that the access policy is assigned to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Access Policy Assignment.
