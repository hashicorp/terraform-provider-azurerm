---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_profile"
description: |-
  Manages a Frontdoor Profile to create a collection of Frontdoor Endpoints.
---

# azurerm_frontdoor_profile

Manages a Frontdoor Profile to create a collection of Frontdoor Endpoints.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "example" {
  name                = "exampleFrontdoorProfile"
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard_Verizon"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Frontdoor Profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Frontdoor Profile.

* `sku_name` - (Required) The pricing related information of current Frontdoor profile. Accepted values are `Custom_Verizon`, `Premium_AzureFrontDoor`, `Premium_Verizon`, `Standard_Akamai`, `Standard_AvgBandWidth_ChinaCdn`, `Standard_AzureFrontDoor`, `Standard_ChinaCdn`, `Standard_Microsoft`, `Standard_955BandWidth_ChinaCdn`, `StandardPlus_AvgBandWidth_ChinaCdn`, `StandardPlus_ChinaCdn`, `StandardPlus_955BandWidth_ChinaCdn` or `Standard_Verizon`.

* `identity` - (Optional) An `identity` block as defined below.

* `origin_response_timeout_seconds` - Possible values are between `16` and `240` seconds(inclusive). Defaults to `120` seconds.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
An `identity` block supports the following:

* `type` - (Required) The type of identity used for the Frontdoor Profile. Possible values are `None`, `UserAssigned`, `SystemAssigned` or `SystemAssigned, UserAssigned`. 

* `identity_ids` - (Optional) A list of User Assigned Identity IDs which should be assigned to this Frontdoor Profile.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Frontdoor Profile.

* `frontdoor_id` - The UUID of the Frontdoor instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Profile.
* `update` - (Defaults to 30 minutes) Used when updating the Frontdoor Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Profile.

## Import

Frontdoor Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1
```
