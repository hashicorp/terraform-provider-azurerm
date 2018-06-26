---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_plan"
sidebar_current: "docs-azurerm-resource-app-service-plan"
description: |-
  Create an App Service Plan component.
---

# azurerm_app_service_plan

Create an App Service Plan component.

## Example Usage (Dedicated)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "api-appserviceplan-pro"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}
```

## Example Usage (Shared / Consumption Plan)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "api-appserviceplan-pro"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "FunctionApp"

  sku {
    tier = "Dynamic"
    size = "Y1"
  }
}
```

## Example Usage (Linux)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "api-appserviceplan-pro"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Linux"

  sku {
    tier = "Standard"
    size = "S1"
  }

  properties {
    reserved = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service Plan component. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service Plan component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `kind` - (Optional) The kind of the App Service Plan to create. Possible values are `Windows` (also available as `App`), `Linux` and `FunctionApp` (for a Consumption Plan). Defaults to `Windows`. Changing this forces a new resource to be created.

~> **NOTE:** When creating a `Linux` App Service Plan, the `reserved` field must be set to `true`.

* `sku` - (Required) A `sku` block as documented below.

* `properties` - (Optional) A `properties` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

`sku` supports the following:

* `tier` - (Required) Specifies the plan's pricing tier.

* `size` - (Required) Specifies the plan's instance size.

* `capacity` - (Optional) Specifies the number of workers associated with this App Service Plan.

`properties` supports the following:

* `app_service_environment_id` - (Optional) The ID of the App Service Environment where the App Service Plan should be located. Changing forces a new resource to be created.

~> **NOTE:** Attaching to an App Service Environment requires the App Service Plan use a `Premium` SKU (when using an ASEv1) and the `Isolated` SKU (for an ASEv2).

* `maximum_number_of_workers` - (Optional) Maximum number of instances that can be assigned to this App Service plan.

* `reserved` - (Optional) Is this App Service Plan `Reserved`. Defaults to `false`.

* `per_site_scaling` - (Optional) Can Apps assigned to this App Service Plan be scaled independently? If set to `false` apps assigned to this plan will scale to all instances of the plan.  Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Plan component.
* `maximum_number_of_workers` - The maximum number of workers supported with the App Service Plan's sku.

## Import

App Service Plan instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_plan.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1
```
