---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_management_private_link_association"
description: |-
  Manages a Resource Management Private Link Association.
---

# azurerm_resource_management_private_link_association

Manages a Resource Management Private Link Association.

## Example Usage

```hcl
data "azurerm_client_config" "example" {}

data "azurerm_management_group" "example" {
  name = data.azurerm_client_config.example.tenant_id
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_resource_management_private_link" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "random_uuid" "example" {
}

resource "azurerm_resource_management_private_link_association" "example" {
  name                                = random_uuid.example.result
  management_group_id                 = azurerm_management_group.example.id
  resource_management_private_link_id = azurerm_resource_management_private_link.example.id
  public_network_access_enabled       = true
}

```

## Arguments Reference

The following arguments are supported:

* `management_group_id` - (Required) Specifies the Management Group ID within which this Private Link Association should exist. Changing this forces a new Private Link Association to be created.

~> **Note:** For now, `management_group_id` must be the ID of [Root Management Group](https://learn.microsoft.com/en-us/azure/governance/management-groups/overview#root-management-group-for-each-directory).

* `resource_management_private_link_id` - (Required) The Resource ID of Resource Management Private Link. Changing this forces a new Private Link Association to be created.

* `public_network_access_enabled` - (Required) Whether public network access is allowed. Changing this forces a new Private Link Association to be created.
 
* `name` - (Optional) Specifies the name of this Private Link Association, which should be a UUID. If `name` is not provided, a UUID will be generated, you should use [the `ignore_changes` attribute to ignore changes to this field](https://www.terraform.io/language/meta-arguments/lifecycle#ignore_changess). Changing this forces a new Private Link Association to be created.

```hcl
resource "azurerm_resource_management_private_link_association" "example" {
  management_group_id                 = azurerm_management_group.example.id
  resource_management_private_link_id = azurerm_resource_management_private_link.example.id
  public_network_access_enabled       = true
  lifecycle {
    ignore_changes = [name]
  }
}

```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Private Link Association.

* `tenant_id` - The Tenant ID.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private Link Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private Link Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private Link Association.

## Import

An existing Private Link Association can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_resource_management_private_link_association.example /providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/privateLinkAssociations/00000000-0000-0000-0000-000000000000
```
