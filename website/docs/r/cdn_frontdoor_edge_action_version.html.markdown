---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_edge_action_version"
description: |-
  Manages a version of a Front Door (standard/premium) Edge Action.
---

# azurerm_cdn_frontdoor_edge_action_version

Manages a version of a Front Door (standard/premium) Edge Action, including the code it runs.

~> **Note:** The first version created for an Edge Action becomes its default version. The default version can only be deleted once all other versions of the Edge Action have been deleted, so additional versions should depend on the default version (for example using `depends_on`) to ensure they are destroyed first.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_edge_action" "example" {
  name                = "exampleedgeaction"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cdn_frontdoor_edge_action_version" "example" {
  name                         = "v1"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.example.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = filebase64("edge_action.js")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Edge Action Version. Possible values must be between 1 and 50 characters in length and may only contain letters and numbers. Changing this forces a new resource to be created.

* `cdn_frontdoor_edge_action_id` - (Required) The ID of the Front Door Edge Action this version belongs to. Changing this forces a new resource to be created.

* `deployment_type` - (Required) The deployment type of the code for this Edge Action Version. Possible values are `file`, `others` and `zip`. Changing this forces a new resource to be created.

---

* `code_content_base64` - (Optional) The base64 encoded content of the code to deploy to this Edge Action Version, either a JavaScript file or a zip archive depending on `deployment_type`.

* `code_file_name` - (Optional) The file name of the code to deploy to this Edge Action Version, for example `edge_action.js`.

-> **Note:** `code_content_base64` and `code_file_name` must be specified together. Deploying the code can take more than 10 minutes, and the code is validated by the service when it is deployed. The deployed code can't be read back from the API in a form which matches the configuration, so these fields aren't populated when a version is imported and changes made to the code outside of Terraform aren't detected.

~> **Note:** The code of a version which isn't the default can be changed in-place. Changing the deployed code of the default version, or removing the code from any version, forces a new resource to be created - however the default version can't be deleted while the Edge Action has other versions, so this plan will fail to apply once other versions exist. To roll out new code, deploy it to a version which isn't the default and then make that version the default (which currently needs to be done outside of Terraform, for example in the Azure Portal).

* `tags` - (Optional) A mapping of tags which should be assigned to the Edge Action Version.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Edge Action Version.

* `default_version_enabled` - Whether this is the default version of the Edge Action.

* `validation_status` - The validation status of the code deployed to this Edge Action Version. This is empty until code has been deployed to the version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Edge Action Version.
* `read` - (Defaults to 5 minutes) Used when retrieving the Edge Action Version.
* `update` - (Defaults to 1 hour) Used when updating the Edge Action Version.
* `delete` - (Defaults to 30 minutes) Used when deleting the Edge Action Version.

## Import

Edge Action Versions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_edge_action_version.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1/versions/version1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2026-10-01
