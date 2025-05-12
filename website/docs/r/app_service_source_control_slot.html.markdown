---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_source_control_slot"
description: |-
  Manages an App Service Source Control Slot.
---

# azurerm_app_service_source_control_slot

Manages an App Service Source Control Slot.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_service_plan" "example" {
  name                = "example-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example-web-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_linux_web_app_slot" "example" {
  name           = "example-slot"
  app_service_id = azurerm_linux_web_app.example.id

  site_config {}
}

resource "azurerm_app_service_source_control_slot" "example" {
  slot_id  = azurerm_linux_web_app_slot.example.id
  repo_url = "https://github.com/Azure-Samples/python-docs-hello-world"
  branch   = "master"
}

```

## Arguments Reference

The following arguments are supported:

* `slot_id` - (Required) The ID of the Linux or Windows Web App Slot. Changing this forces a new resource to be created.

~> **Note:** Function App Slots are not supported at this time.

---

* `branch` - (Optional) The URL for the repository. Changing this forces a new resource to be created.

* `github_action_configuration` - (Optional) A `github_action_configuration` block as detailed below. Changing this forces a new resource to be created.

* `repo_url` - (Optional) The branch name to use for deployments. Changing this forces a new resource to be created.

* `rollback_enabled` - (Optional) Should the Deployment Rollback be enabled? Defaults to `false` Changing this forces a new resource to be created.

* `use_local_git` - (Optional) Should the Slot use local Git configuration. Changing this forces a new resource to be created.

* `use_manual_integration` - (Optional) Should code be deployed manually. Set to `true` to disable continuous integration, such as webhooks into online repos such as GitHub. Defaults to `false`. Changing this forces a new resource to be created.

* `use_mercurial` - (Optional) The repository specified is Mercurial. Defaults to `false`. Changing this forces a new resource to be created.

---

A `github_action_configuration` block supports the following:

* `code_configuration` - (Optional) A `code_configuration` block as detailed below. Changing this forces a new resource to be created.

* `container_configuration` - (Optional) A `container_configuration` block as detailed below.

* `generate_workflow_file` - (Optional) Should the service generate the GitHub Action Workflow file. Defaults to `true` Changing this forces a new resource to be created.

* `linux_action` - Denotes this action uses a Linux base image.

---

A `code_configuration` block supports the following:

* `runtime_stack` - (Required) The value to use for the Runtime Stack in the workflow file content for code base apps. Changing this forces a new resource to be created. Possible values are `dotnetcore`, `spring`, `tomcat`, `node` and `python`.

* `runtime_version` - (Required) The value to use for the Runtime Version in the workflow file content for code base apps. Changing this forces a new resource to be created.

---

A `container_configuration` block supports the following:

* `image_name` - (Required) The image name for the build. Changing this forces a new resource to be created.

* `registry_password` - (Optional) The password used to upload the image to the container registry. Changing this forces a new resource to be created.

* `registry_url` - (Required) The server URL for the container registry where the build will be hosted. Changing this forces a new resource to be created.

* `registry_username` - (Optional) The username used to upload the image to the container registry. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Source Control Slot

* `scm_type` - The SCM Type in use. This value is decoded by the service from the repository information supplied.

* `uses_github_action` - Indicates if the Slot uses a GitHub action for deployment. This value is decoded by the service from the repository information supplied.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Source Control Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source Control Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Source Control Slot.

## Import

an App Service Source Control Slot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_source_control_slot.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1"
```
