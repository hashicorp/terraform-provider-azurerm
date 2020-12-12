---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_ssh_public_key"
description: |-
  Gets information about an existing SSH Public Key.
---

# Data Source: azurerm_ssh_public_key

Use this data source to access information about an existing SSH Public Key.

## Example Usage

```hcl
data "azurerm_ssh_public_key" "example" {
  name = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_ssh_public_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this SSH Public Key.

* `resource_group_name` - (Required) The name of the Resource Group where the SSH Public Key exists.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the SSH Public Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the SSH Public Key.

* `public_key` - The SSH public key used to authenticate to a virtual machine through ssh.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SSH Public Key.