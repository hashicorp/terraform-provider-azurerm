---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Manages a Front Door (standard/premium) Profile.
---

# azurerm_cdn_frontdoor_profile

Manages a Front Door (standard/premium) Profile which contains a collection of endpoints and origin groups.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-identity"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                     = "example-cdn-profile"
  resource_group_name      = azurerm_resource_group.example.name
  sku_name                 = "Premium_AzureFrontDoor"
  response_timeout_seconds = 120

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  log_scrubbing_rule {
    match_variable = "RequestIPAddress"
  }

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this Front Door Profile should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU for this Front Door Profile. Possible values include `Standard_AzureFrontDoor` and `Premium_AzureFrontDoor`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `response_timeout_seconds` - (Optional) Specifies the maximum response timeout in seconds. Possible values are between `16` and `240` seconds (inclusive). Defaults to `120` seconds.

* `log_scrubbing_rule` - (Optional) One or more `log_scrubbing_rule` blocks as defined below.

~> **Note:** When no `log_scrubbing_rule` blocks are defined, log scrubbing will be automatically `disabled`. When one or more `log_scrubbing_rule` blocks are present, log scrubbing will be `enabled`.

* `tags` - (Optional) Specifies a mapping of tags to assign to the resource.


---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible values are `SystemAssigned`, `UserAssigned` or `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `log_scrubbing_rule` block supports the following:

* `match_variable` - (Required) The variable to be scrubbed from the logs. Possible values are `QueryStringArgNames`, `RequestIPAddress`, and `RequestUri`.

~> **Note:** The `operator` field is implicitly set to `EqualsAny`, as it is the sole supported value, and is therefore not exposed as a configurable option in the provider schema.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this Front Door Profile.

* `resource_guid` - The UUID of this Front Door Profile which will be sent in the HTTP Header as the `X-Azure-FDID` attribute.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Front Door Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Front Door Profile.
* `delete` - (Defaults to 2 hours) Used when deleting the Front Door Profile.

## Import

Front Door Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2024-02-01
