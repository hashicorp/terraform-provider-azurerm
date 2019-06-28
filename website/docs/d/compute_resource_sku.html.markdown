---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_compute_resource_sku"
sidebar_current: "docs-azurerm-datasource-compute-resource-sku"
description: |-
  Retrieve information about a Compute SKU

---

# Data Source: azurerm_compute_resource_sku

Use this data source to access information about Compute Resource SKUs.

## Example Usage

```hcl
data "azurerm_compute_resource_sku" "test" {
  name     = "Standard_D2s_v2"
  location = "eastus"
}

output "location_info" {
  value = "${data.azurerm_compute_resource_sku.test.location_info}"
}

output "zones" {
  value = "${data.azurerm_compute_resource_sku.test.location_info.0.zones}"
}
```

## Argument Reference

* `name` - (Required) The name of the SKU to retrieve information for (e.g. Standard_DS2_v2).
* `location` - (Required) The location to retrieve SKU information from.

## Attributes Reference

The following attributes are exported (Not all attributes may be available for every SKU):

* `resource_type` - The type of resource the SKU applies to.
  
* `name` - The name of the SKU.

* `tier` - Specifies the tier of the virtual machines in a scale set. Possible values: **Standard** or **Basic**.

* `size` - The size of the SKU.

* `family` - The Family of this particular SKU.

* `kind` - The Kind of resources that are supported in this SKU.

* `capacity` - Specifies the number of virtual machines in the scale set.

* `location_info` - A list of locations and availability zones in those locations where the SKU is available.

* `api_versions` - The api versions that support this SKU.

* `costs` - Metadata for retrieving price info.

* `capabilities` - A set of name value pairs describing capabilities.

* `restrictions` - The restrictions because of which SKU cannot be used. This is empty if there are no restrictions.
  