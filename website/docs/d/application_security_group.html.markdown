---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_security_group"
sidebar_current: "docs-azurerm-datasource-network-application-security-group"
description: |-
  Get information about an Application Security Group.
---

# Data Source: azurerm_application_security_group

Get information about an Application Security Group.

## Example Usage

```hcl
data "azurerm_application_security_group" "test" {
  name = "tf-appsecuritygroup"
  resource_group_name = "my-resource-group"
}

output "application_security_group_id" {
  value = "${data.azurerm_application_security_group.test.id}"
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
