---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ssh_public_key"
description: |-
  Manages a SSH Public Key.
---

# azurerm_ssh_public_key

Manages a SSH Public Key.

## Example Usage

```hcl
resource "azurerm_ssh_public_key" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
  public_key          = file("~/.ssh/id_rsa.pub")
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the SSH Public Key should exist. Changing this forces a new SSH Public Key to be created.

* `name` - (Required) The name which should be used for this SSH Public Key. Changing this forces a new SSH Public Key to be created.

* `public_key` - (Required) SSH public key used to authenticate to a virtual machine through ssh. the provided public key needs to be at least 2048-bit and in ssh-rsa format.

* `resource_group_name` - (Required) The name of the Resource Group where the SSH Public Key should exist. Changing this forces a new SSH Public Key to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the SSH Public Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the SSH Public Key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the SSH Public Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the SSH Public Key.
* `update` - (Defaults to 45 minutes) Used when updating the SSH Public Key.
* `delete` - (Defaults to 45 minutes) Used when deleting the SSH Public Key.

## Import

SSH Public Keys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_ssh_public_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/SshPublicKeys/mySshPublicKeyName1
```
