---
subcategory: "Arc Resource Bridge"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_arc_resource_bridge_appliance"
description: |-
  Gets information about an existing Arc Resource Bridge Appliance.
---

# Data Source: azurerm_arc_resource_bridge_appliance

Use this data source to access information about an existing Arc Resource Bridge Appliance.

## Example Usage

```hcl
data "azurerm_arc_resource_bridge_appliance" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_arc_resource_bridge_appliance.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Arc Resource Bridge Appliance.

* `resource_group_name` - (Required) The name of the Resource Group where the Arc Resource Bridge Appliance exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Arc Resource Bridge Appliance.

* `distro` - Fabric/Infrastructure for this Arc Resource Bridge Appliance.

* `identity` - An `identity` block as defined below.

* `infrastructure_provider` - The infrastructure provider about the connected Arc Resource Bridge Appliance.

* `location` - The Azure Region where the Arc Resource Bridge Appliance exists.

* `public_key_base64` - RSA public key in PKCS1 format encoded in base64.

* `tags` - A mapping of tags assigned to the Arc Resource Bridge Appliance.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The type of this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Resource Bridge Appliance.
