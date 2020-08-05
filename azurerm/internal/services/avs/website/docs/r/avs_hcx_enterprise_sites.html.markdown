---
subcategory: "Avs"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_avs_hcx_enterprise_site"
description: |-
  Manages a avs HcxEnterpriseSite.
---

# azurerm_avs_hcx_enterprise_site

Manages a avs HcxEnterpriseSite.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_avs_private_cloud" "example" {
  name = "example-privatecloud"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
  sku {
      name = "example-privatecloud"
  }

  management_cluster {
      cluster_size = 42
  }
  network_block = ""
}

resource "azurerm_avs_hcx_enterprise_site" "example" {
  name = "example-hcxenterprisesite"
  resource_group_name = azurerm_resource_group.example.name
  private_cloud_name = azurerm_avs_private_cloud.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this avs HcxEnterpriseSite. Changing this forces a new avs HcxEnterpriseSite to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the avs HcxEnterpriseSite should exist. Changing this forces a new avs HcxEnterpriseSite to be created.

* `private_cloud_name` - (Required) The name of the private cloud. Changing this forces a new avs HcxEnterpriseSite to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the avs HcxEnterpriseSite.

* `activation_key` - The activation key.

* `type` - Resource type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the avs HcxEnterpriseSite.
* `read` - (Defaults to 5 minutes) Used when retrieving the avs HcxEnterpriseSite.
* `delete` - (Defaults to 30 minutes) Used when deleting the avs HcxEnterpriseSite.

## Import

avs HcxEnterpriseSites can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_avs_hcx_enterprise_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AVS/privateClouds/privateCloud1/hcxEnterpriseSites/hcxEnterpriseSite1
```