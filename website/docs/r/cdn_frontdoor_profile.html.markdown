---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Manages a Frontdoor Profile to create a collection of Frontdoor Endpoints.
---

# azurerm_cdn_frontdoor_profile

Manages a Frontdoor Profile to create a collection of Frontdoor Endpoints.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-cdn-frontdoor"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_AzureFrontDoor"

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Frontdoor Profile. Possible values must be between 2 and 90 characters in length, begin with a letter or number, end with a letter or number and may contain only letters, numbers and hyphens. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Frontdoor Profile.

* `sku_name` - (Optional) The pricing related information of current Frontdoor profile. Possible values include `Premium_AzureFrontDoor` and `Standard_AzureFrontDoor`. Defaults to `Standard_AzureFrontDoor`.

* `response_timeout_seconds` - Possible values are between `16` and `240` seconds(inclusive). Defaults to `120` seconds.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Frontdoor Profile.

* `cdn_frontdoor_id` - The UUID of the Frontdoor instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Profile.

## Import

Frontdoor Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1
```
