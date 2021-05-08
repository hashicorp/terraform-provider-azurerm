---
subcategory: "Portal"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_tenant_configuration"
description: |-
  Manages Portal Tenant Configuration.
---

# azurerm_tenant_configuration

Manages Portal Tenant Configuration.

~> **Note:** User has to be a Tenant Admin for managing this resource.

## Example Usage

```hcl
resource "azurerm_tenant_configuration" "example" {
  enforce_private_markdown_storage = true
}
```

## Arguments Reference

The following arguments are supported:

* `enforce_private_markdown_storage` - (Required) Is Markdown tile which used to display custom and static content enabled?

~> **Note:** When `enforce_private_markdown_storage` is set to `true` Markdown tile will require external storage configuration (URI). The inline content configuration will be prohibited.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Portal Tenant Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Portal Tenant Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Portal Tenant Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Portal Tenant Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Portal Tenant Configuration.

## Import

Portal Tenant Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_tenant_configuration.example /providers/Microsoft.Portal/tenantConfigurations/default
```
