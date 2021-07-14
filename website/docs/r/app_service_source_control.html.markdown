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
  location            = "West Europe"
  os_type             = "Linux"
  sku_name            = "P1V2"
}

resource "azurerm_linux_web_app" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id
}

resource "azurerm_app_service_source_control" "example" {
  app_id   = azurerm_linux_web_app.example.id
  repo_url = "https://github.com/Azure-Samples/python-docs-hello-world"
  branch   = "master"
}
```

## Arguments Reference

The following arguments are supported:

* `app_id` - (Required) The ID of the Windows or Linux Web App.

* `branch` - (Required) The branch name to use for deployments.

* `repo_url` - (Required) The URL for the repository.

---

* `github_action_configuration` - (Optional) A `github_action_configuration` block as defined below.

* `manual_integration` - (Optional) Should code be deployed manually. Set to `false` to enable continuous integration, such as webhooks into online repos such as GitHub.

* `rollback_enabled` - (Optional) Should the Deployment Rollback be enabled? Defaults to `false`

* `scm_type` - (Optional) The SCM System to use for Source Control. Possible values include 'ScmTypeNone', 'ScmTypeDropbox', 'ScmTypeTfs', 'ScmTypeLocalGit', 'ScmTypeGitHub', 'ScmTypeCodePlexGit', 'ScmTypeCodePlexHg', 'ScmTypeBitbucketGit', 'ScmTypeBitbucketHg', 'ScmTypeExternalGit', 'ScmTypeExternalHg', 'ScmTypeOneDrive', 'ScmTypeVSO'.

~> **NOTE:** Azure can typically set this value automatically based on the `repo_url` value. 

~> **NOTE:** SCM Type `ScmTypeVSTSRM` is not supported as this is set by Azure DevOps and overrides Terraform's control of this resource.

* `use_mercurial` - (Optional) The repository specified is Mercurial. Defaults to `false`.

* `uses_github_action` - (Optional) Should deployment be performed by GitHub Action. Defaults to `false`.

---

A `code_configuration` block supports the following:

* `runtime_stack` - (Required) The value to use for the Runtime Stack in the workflow file content for code base apps.

* `runtime_version` - (Optional) The value to use for the Runtime Version in the workflow file content for code base apps.

---

A `container_configuration` block supports the following:

* `image_name` - (Required) The image name for the build.

* `registry_url` - (Required) The server URL for the container registry where the build will be hosted.

* `registry_password` - (Optional) The password used to upload the image to the container registry.

* `registry_username` - (Optional) The username used to upload the image to the container registry.

---

A `github_action_configuration` block supports the following:

* `code_configuration` - (Optional) A `code_configuration` block as defined above.

* `container_configuration` - (Optional) A `container_configuration` block as defined above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the App Service Source Control.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Source Control.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Source Control.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Source Control.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Source Control.

## Import

App Service Source Controls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_source_control.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
```