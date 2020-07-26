---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_service_tags"
description: |-
  Gets information about Service Tags for a specific service type.
---

# Data Source: azurerm_network_service_tags

Use this data source to access information about Service Tags.

## Example Usage

```hcl
data "azurerm_network_service_tags" "example" {
  location        = "westcentralus"
  service         = "AzureKeyVault"
  location_filter = "northeurope"
}

output "address_prefixes" {
  value = data.azurerm_service_tags.example.address_prefixes
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Service Tags exists. This value is not used to filter the results but for specifying the region to request. For filtering by region use `location_filter` instead.  More information can be found here: [Service Tags URL parameters](https://docs.microsoft.com/en-us/rest/api/virtualnetwork/servicetags/list#uri-parameters).

* `service` - (Required) The type of the service for which address prefixes will be fetched. Available service tags can be found here: [Available service tags](https://docs.microsoft.com/en-us/azure/virtual-network/service-tags-overview#available-service-tags).

---

* `location_filter` - (Optional) Changes the scope of the service tags. Can be any value that is also valid for `location`. If this field is empty then all address prefixes are considered instead of only location specific ones.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of this Service Tags block.

* `address_prefixes` - List of address prefixes for the service type (and optionally a specific region).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Service Tags.
