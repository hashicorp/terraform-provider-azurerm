---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_source_control"
description: |-
  Manages source control for an App Service.

---

# azurerm_app_service_source_control

Manages source control for an App Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-app-service-plan"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app-service"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  app_service_plan_id = "${azurerm_app_service_plan.example.id}"

  lifecycle {
    ignore_changes = [site_config.0.scm_type]
  }
}

resource "azurerm_app_service_source_control" "example" {
  app_service_id             = "${azurerm_app_service.example.id}"
  repo_url                   = "https://github.com/Azure-Samples/app-service-web-html-get-started"
  branch                     = "master"
  manual_integration_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `repo_url` - (Required) The repository or source control URL. Changing this forces a new resource to be created.

* `branch` - (Optional) The name of branch to use for deployment. Changing this forces a new resource to be created.

* `deployment_rollback_enabled` - (Optional) Should deployment rollback be enabled? Defaults to `false`. Changing this forces a new resource to be created.

* `is_manual_integration` - (Optional) Should manual integration be enabled, rather than continuous integration (which configures webhooks into online repos like GitHub)? Defaults to `false`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Source Control.

## Import

App Service source control can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_source_control.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-group/providers/Microsoft.Web/sites/test-app/sourcecontrols/web
```
