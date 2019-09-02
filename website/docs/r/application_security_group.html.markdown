---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_security_group"
sidebar_current: "docs-azurerm-resource-network-application-security-group"
description: |-
  Manages an Application Security Group.
---

# azurerm_application_security_group

Manages an Application Security Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_security_group" "test" {
  name                = "tf-appsecuritygroup"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    Hello = "World"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Security Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Application Security Group.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Application Security Group.

## Import

Application Security Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_security_group.securitygroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/applicationSecurityGroups/securitygroup1
```
