---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_profile"
description: |-
  Gets information about an existing Front Door (standard/premium) Profile.
---

# Data Source: azurerm_cdn_frontdoor_profile

Use this data source to access information about an existing Front Door (standard/premium) Profile.

## Example Usage

```hcl
data "azurerm_cdn_frontdoor_profile" "example" {
  name                = "existing-cdn-profile"
  resource_group_name = "existing-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Front Door Profile.

* `resource_group_name` - (Required) The name of the Resource Group where this Front Door Profile exists.

* `identity` - (Optional) An `identity` block as defined below.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Front Door Profile.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Front Door Profile.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of this Front Door Profile.

* `resource_guid` - The UUID of the Front Door Profile which will be sent in the HTTP Header as the `X-Azure-FDID` attribute.

* `sku_name` - Specifies the SKU for this Front Door Profile.

* `response_timeout_seconds` - Specifies the maximum response timeout in seconds.

* `tags` - Specifies a mapping of Tags assigned to this Front Door Profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Front Door Profile.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Cdn`: 2024-02-01
