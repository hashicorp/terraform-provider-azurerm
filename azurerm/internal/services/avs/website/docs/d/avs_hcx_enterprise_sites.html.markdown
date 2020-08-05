---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_avs_hcx_enterprise_site"
description: |-
  Gets information about an existing avs HcxEnterpriseSite.
---

# Data Source: azurerm_avs_hcx_enterprise_site

Use this data source to access information about an existing avs HcxEnterpriseSite.

## Example Usage

```hcl
data "azurerm_avs_hcx_enterprise_site" "example" {
  name = "example-hcxenterprisesite"
  resource_group_name = "example-resource-group"
  private_cloud_name = "existing"
}

output "id" {
  value = data.azurerm_avs_hcx_enterprise_site.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this avs HcxEnterpriseSite.

* `resource_group_name` - (Required) The name of the Resource Group where the avs HcxEnterpriseSite exists.

* `private_cloud_name` - (Required) Name of the private cloud.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs HcxEnterpriseSite.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the avs HcxEnterpriseSite.