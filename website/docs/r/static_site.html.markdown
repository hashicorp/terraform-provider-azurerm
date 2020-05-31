---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_site"
description: |-
  Manages a Static Web App.
---

# azurerm_static_site

Manages a Static Web App.

## Example Usage

```hcl
resource "azurerm_static_site" "example" {
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"

  github_configuration {
    repo_token   = "personal-access-token-github"
    repo_url     = "https://github.com/example/static-web-app-example"
    branch       = "master"
    app_location = "/"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Static Web App. Changing this forces a new Static Web App to be created.

* `location` - (Required) The Azure Region where the Static Web App should exist. Changing this forces a new Static Web App to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Static Web App should exist. Changing this forces a new Static Web App to be created.

* `github_configuration` - (Required) A `github_configuration` block as defined below.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Static Web App.

---

A `github_configuration` block supports the following:

* `app_location` - (Required) The path to the Static Web App site code within the repository.

* `branch` - (Required) The target branch in the repository.

* `repo_token` - (Required) A user's github repository token. This is used to setup the Github Actions workflow file and API secrets.

* `repo_url` - (Required) URL for the repository of the Static Web App site.

* `api_location` - (Optional) The path to the Function App api code within the repository.

* `artifact_location` - (Optional) The path of the Static Web App artifacts after building.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Static Web App.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Web App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App.
* `update` - (Defaults to 30 minutes) Used when updating the Static Web App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Web App.

## Import

Static Web Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1
```
