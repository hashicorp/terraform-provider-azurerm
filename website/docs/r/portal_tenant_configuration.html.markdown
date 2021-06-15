---
subcategory: "Portal"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_portal_tenant_configuration"
description: |-
  Manages Portal Tenant Configuration.
---

# azurerm_portal_tenant_configuration

Manages Portal Tenant Configuration.

~> **Note:** User has to be a Tenant Admin for managing this resource. The Service Principal with Tenant Admin could be created by `az ad sp create-for-rbac --name "<sp name>" --role="Owner" --scopes="/"`.

## Example Usage

```hcl
resource "azurerm_portal_tenant_configuration" "example" {
  private_markdown_storage_enforced = true
}
```

## Arguments Reference

The following arguments are supported:

* `private_markdown_storage_enforced` - (Required) Is the private tile markdown storage which used to display custom dynamic and static content enabled?

~> **Note:** When `private_markdown_storage_enforced` is set to `true` Markdown tile will require external storage configuration (URI). The inline content configuration will be prohibited.

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
terraform import azurerm_portal_tenant_configuration.example /providers/Microsoft.Portal/tenantConfigurations/default
```
