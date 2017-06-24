---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_plan"
sidebar_current: "docs-azurerm-resource-app_service_plan"
description: |-
  Create an App Service Plan component.
---

# azurerm\_app\_service\_plan

Create an App Service Plan component.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "api-appserviceplan-pro"
  location            = "West Europe"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
		name = "Standard_S0"
		tier = "Standard"
		size = "S0"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service Plan component. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the App Service Plan component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

`sku` supports the following:

* `tier` - (Required) .
* `size` - (Required) .

* `maximum_number_of_workers` - (Optional) .

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Plan component.

## Import

App Service Plan instances can be imported using the `resource id`, e.g.

```
terraform import azurerm_app_service_plan.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1
```
