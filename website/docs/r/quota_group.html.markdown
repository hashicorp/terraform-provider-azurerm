---
subcategory: "Quota"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_quota_group"
description: |-
  Manages an Azure Quota Group.
---

# azurerm_quota_group

Manages an Azure Quota Group, which allows organizations to pool and centrally manage virtual machine quota across multiple subscriptions under a Management Group.

~> **Note:** Quota Groups are a feature of the [Azure Quotas service](https://learn.microsoft.com/en-us/azure/quotas/quota-groups) and require the `Microsoft.Quota` resource provider to be registered on the Management Group's subscriptions.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "example" {
  name                = "example-quota-group"
  management_group_id = data.azurerm_management_group.example.id
  display_name        = "Example Quota Group"
}
```

## Example Usage with Subscription Associations and Quota Requests

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_quota_group" "example" {
  name                = "example-quota-group"
  management_group_id = data.azurerm_management_group.example.id
  display_name        = "Example Quota Group"

  associated_subscription_ids = [
    "00000000-0000-0000-0000-000000000001",
    "00000000-0000-0000-0000-000000000002",
  ]

  quota_request {
    resource_name = "standardddv4family"
    location      = "eastus"
    limit         = 50
    comment       = "Pooled quota for DDv4 workloads"
  }

  quota_request {
    resource_name = "standardDSv3Family"
    location      = "eastus"
    limit         = 100
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Quota Group. Changing this forces a new Quota Group to be created.

* `management_group_id` - (Required) The ID of the Management Group this Quota Group is scoped to. Changing this forces a new Quota Group to be created.

---

* `display_name` - (Optional) A human-readable display name for this Quota Group.

* `associated_subscription_ids` - (Optional) A set of Subscription UUIDs to associate with this Quota Group.

* `quota_request` - (Optional) One or more `quota_request` blocks as defined below.

---

A `quota_request` block supports the following:

* `resource_name` - (Required) The name of the quota resource (e.g. `standardDDv4Family`).

* `location` - (Required) The Azure region for this quota request (e.g. `eastus`).

* `limit` - (Required) The requested quota limit for the resource in the given location.

* `resource_provider_name` - (Optional) The resource provider namespace. Defaults to `Microsoft.Compute`.

* `comment` - (Optional) A comment explaining the purpose of this quota request.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Quota Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Quota Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Quota Group.
* `update` - (Defaults to 30 minutes) Used when updating the Quota Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Quota Group.

## Import

Quota Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_quota_group.example /providers/Microsoft.Management/managementGroups/example-management-group/providers/Microsoft.Quota/groupQuotas/example-quota-group
```
