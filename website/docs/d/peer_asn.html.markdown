---
subcategory: "Peering"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_peer_asn"
description: |-
  Gets information about an existing Peer ASN.
---

# Data Source: azurerm_peer_asn

Use this data source to access information about an existing Peer ASN.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_peer_asn" "example" {
  name = "existing"
}

output "id" {
  value = data.azurerm_peer_asn.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Peer ASN.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Peer ASN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Peer ASN.
