---
subcategory: "Peering"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_peer_asn"
description: |-
  Manages a Peer ASN.
---

# azurerm_peer_asn

Manages a Peer ASN.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_peer_asn" "example" {
  name = "example_peerasn"
  asn  = 123
  contact {
    role  = "Noc"
    email = "example@email.com"
  }
  peer_name = "example-peer"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Peer ASN. Changing this forces a new Peer ASN to be created.

* `asn` - (Required) The public Atonomous System Number. Changing this forces a new Peer ASN to be created.

* `contact` - (Required) One or more `contact` blocks as defined below.

* `peer_name` - (Required) The name of the Peer, which needs to be as close as possible to your PeeringDB profile.

---

A `contact` block supports the following:

* `role` - (Required) The role of the contact. Possible values are "Noc", "Other", "Policy", "Service" and "Technical".

* `email` - (Required) The e-mail address of the contact.

* `phone` - (Optional) The phone number of the contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Peer ASN.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Peer ASN.
* `read` - (Defaults to 5 minutes) Used when retrieving the Peer ASN.
* `update` - (Defaults to 30 minutes) Used when updating the Peer ASN.
* `delete` - (Defaults to 30 minutes) Used when deleting the Peer ASN.

## Import

Peer ASNs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_peer_asn.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Peering/peerAsns/peerasn1
```
