---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_profile"
description: |-
  Manages a CDN Profile to create a collection of CDN Endpoints.
---

# azurerm_cdn_profile

Manages a CDN Profile to create a collection of CDN Endpoints.

!> **Note:** Azure rolled out a breaking change on Friday 9th April 2021 which may cause issues with the CDN/FrontDoor resources. [More information is available in this GitHub issue](https://github.com/hashicorp/terraform-provider-azurerm/issues/11231) - unfortunately this may necessitate a breaking change to the CDN and FrontDoor resources, more information will be posted [in the GitHub issue](https://github.com/hashicorp/terraform-provider-azurerm/issues/11231) as the necessary changes are identified.

!> **Note:** The CDN services from Edgio(formerly Verizon) was shut down on 15 January 2025 and is no longer available.

!> **Note:** Support for CDN services from Akamai was removed on 31 October 2023.
## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_profile" "example" {
  name                = "exampleCdnProfile"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard_Microsoft"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the CDN Profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the CDN Profile. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) The pricing related information of current CDN profile. Accepted values are `Standard_Akamai`, `Standard_ChinaCdn`, `Standard_Microsoft`, `Standard_Verizon` or `Premium_Verizon`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the CDN Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CDN Profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the CDN Profile.
* `update` - (Defaults to 30 minutes) Used when updating the CDN Profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the CDN Profile.

## Import

CDN Profiles can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Cdn/profiles/myprofile1
```
