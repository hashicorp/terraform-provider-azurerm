---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_bastion_host"
description: |-
  Gets information about an existing Bastion Host.

---

# Data Source: azurerm_bastion_host

Use this data source to access information about an existing Bastion Host.

## Example Usage

```hcl
data "azurerm_bastion_host" "example" {
  name                = "existing-bastion"
  resource_group_name = "existing-resources"
}

output "id" {
  value = data.azurerm_bastion_host.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Bastion Host.

* `resource_group_name` - (Required) The name of the Resource Group where the Bastion Host exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Bastion Host.

* `location` - The Azure Region where the Bastion Host exists.

* `copy_paste_enabled` - Is Copy/Paste feature enabled for the Bastion Host.

* `file_copy_enabled` - Is File Copy feature enabled for the Bastion Host.

* `sku` - The SKU of the Bastion Host.

* `ip_configuration` - A `ip_configuration` block as defined below.

* `ip_connect_enabled` - Is IP Connect feature enabled for the Bastion Host.

* `scale_units` - The number of scale units provisioned for the Bastion Host.

* `shareable_link_enabled` - Is Shareable Link feature enabled for the Bastion Host.

* `tunneling_enabled` - Is Tunneling feature enabled for the Bastion Host.

* `session_recording_enabled` - Is Session Recording feature enabled for the Bastion Host.

* `dns_name` - The FQDN for the Bastion Host.

* `tags` - A mapping of tags assigned to the Bastion Host.

---

A `ip_configuration` block supports the following:

* `name` - The name of the IP configuration.

* `subnet_id` - Reference to the subnet in which this Bastion Host has been created.

* `public_ip_address_id` - Reference to a Public IP Address associated to this Bastion Host.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Bastion Host.
