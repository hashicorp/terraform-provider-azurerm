---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_security_group"
description: |-
  Gets information about an existing Application Security Group.
---

# Data Source: azurerm_application_security_group

Use this data source to access information about an existing Application Security Group.

## Example Usage

```hcl
data "azurerm_application_security_group" "example" {
  name                = "tf-appsecuritygroup"
  resource_group_name = "my-resource-group"
}

output "application_security_group_id" {
  value = data.azurerm_application_security_group.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the Application Security Group.

* `resource_group_name` - The name of the resource group in which the Application Security Group exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Security Group.

* `location` - The supported Azure location where the Application Security Group exists.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Security Group.
