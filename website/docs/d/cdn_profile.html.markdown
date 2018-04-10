---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_profile"
sidebar_current: "docs-azurerm-datasource-cdn-profile"
description: |-
  Gets information about a CDN Profile
---

# Data Source: azurerm_cdn_profile

Use this data source to access information about a CDN Profile.

## Example Usage

```hcl
data "azurerm_cdn_profile" "test" {
  name = "myfirstcdnprofile"
  resource_group_name = "example-resources"
}

output "cdn_profile_id" {
  value = "${data.azurerm_cdn_profile.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the CDN Profile.

* `resource_group_name` - (Required) The name of the resource group in which the CDN Profile exists.

## Attributes Reference

* `location` - The Azure Region where the resource exists.

* `sku` - The pricing related information of current CDN profile.

* `tags` - A mapping of tags assigned to the resource.
