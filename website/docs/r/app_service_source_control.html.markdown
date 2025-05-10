---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_source_control"
description: |-
  Manages an App Service Web App or Function App Source Control Configuration.
---

# azurerm_app_service_source_control

Manages an App Service Web App or Function App Source Control Configuration.

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
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "P1v2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_app_service_source_control" "example" {
  app_id   = azurerm_linux_web_app.example.id
  repo_url = "https://github.com/Azure-Samples/python-docs-hello-world"
  branch   = "master"
}
```

## Arguments Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the Windows or Linux Web App. Changing this forces a new resource to be created.

~> **Note:** Function apps are not supported at this time.

* `branch` - (Optional) The branch name to use for deployments. Changing this forces a new resource to be created.

* `repo_url` - (Optional) The URL for the repository. Changing this forces a new resource to be created.

---

* `github_action_configuration` - (Optional) A `github_action_configuration` block as defined below. Changing this forces a new resource to be created.

* `use_manual_integration` - (Optional) Should code be deployed manually. Set to `false` to enable continuous integration, such as webhooks into online repos such as GitHub. Defaults to `false`. Changing this forces a new resource to be created.

* `rollback_enabled` - (Optional) Should the Deployment Rollback be enabled? Defaults to `false`. Changing this forces a new resource to be created.

~> **Note:** Azure can typically set this value automatically based on the `repo_url` value.

* `use_local_git` - (Optional) Should the App use local Git configuration. Changing this forces a new resource to be created.

* `use_mercurial` - (Optional) The repository specified is Mercurial. Defaults to `false`. Changing this forces a new resource to be created.

---

A `code_configuration` block supports the following:

* `runtime_stack` - (Required) The value to use for the Runtime Stack in the workflow file content for code base apps. Possible values are `dotnetcore`, `spring`, `tomcat`, `node` and `python`. Changing this forces a new resource to be created.

* `runtime_version` - (Required) The value to use for the Runtime Version in the workflow file content for code base apps. Changing this forces a new resource to be created.

---

A `container_configuration` block supports the following:

* `image_name` - (Required) The image name for the build. Changing this forces a new resource to be created.

* `registry_url` - (Required) The server URL for the container registry where the build will be hosted. Changing this forces a new resource to be created.

* `registry_password` - (Optional) The password used to upload the image to the container registry. Changing this forces a new resource to be created.

* `registry_username` - (Optional) The username used to upload the image to the container registry. Changing this forces a new resource to be created.

---

A `github_action_configuration` block supports the following:

* `code_configuration` - (Optional) A `code_configuration` block as defined above. Changing this forces a new resource to be created.

* `container_configuration` - (Optional) A `container_configuration` block as defined above.

* `generate_workflow_file` - (Optional) Whether to generate the GitHub work flow file. Defaults to `true`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Source Control.

* `uses_github_action` - Indicates if the Slot uses a GitHub action for deployment. This value is decoded by the service from the repository information supplied.

* `scm_type` - The SCM Type in use. This value is decoded by the service from the repository information supplied.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Source Control.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source Control.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Source Control.

## Import

App Service Source Controls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_source_control.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
```
