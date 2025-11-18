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
  object_id           = "00000000-0000-0000-0000-000000000000"
  managed_redis_name  = "example-managedredis"
  resource_group_name = "example-resources"
}

output "object_id" {
  value = data.azurerm_managed_redis_access_policy_assignment.example.object_id
}
```

## Arguments Reference

The following arguments are supported:

* `object_id` - (Required) The object ID of the Azure Active Directory user, group, service principal, or managed identity.

* `managed_redis_name` - (Required) The name of the Managed Redis instance.

* `database_name` - (Optional) The name of the database within the Managed Redis instance. Defaults to `default`.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Redis instance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Redis Access Policy Assignment.

* `object_id` - The object ID of the Azure Active Directory user, group, service principal, or managed identity that the access policy is assigned to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Redis Access Policy Assignment.
