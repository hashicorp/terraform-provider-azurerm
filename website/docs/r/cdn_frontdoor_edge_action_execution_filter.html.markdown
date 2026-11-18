---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_edge_action_execution_filter"
description: |-
  Manages an execution filter for a Front Door (standard/premium) Edge Action.
---

# azurerm_cdn_frontdoor_edge_action_execution_filter

Manages an execution filter for a Front Door (standard/premium) Edge Action. Requests which contain the specified header name and value run the specified Edge Action Version instead of the default version.

~> **Note:** An Edge Action Version which is referenced by an execution filter cannot be deleted until the execution filter has been deleted.

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

resource "azurerm_cdn_frontdoor_edge_action_version" "v1" {
  name                         = "v1"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.example.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = filebase64("v1/edge_action.js")
}

resource "azurerm_cdn_frontdoor_edge_action_version" "v2" {
  name                         = "v2"
  cdn_frontdoor_edge_action_id = azurerm_cdn_frontdoor_edge_action.example.id
  deployment_type              = "file"
  code_file_name               = "edge_action.js"
  code_content_base64          = filebase64("v2/edge_action.js")

  depends_on = [azurerm_cdn_frontdoor_edge_action_version.v1]
}

resource "azurerm_cdn_frontdoor_edge_action_execution_filter" "example" {
  name                                 = "canary"
  cdn_frontdoor_edge_action_id         = azurerm_cdn_frontdoor_edge_action.example.id
  cdn_frontdoor_edge_action_version_id = azurerm_cdn_frontdoor_edge_action_version.v2.id
  header_name                          = "x-canary"
  header_value                         = "true"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Edge Action Execution Filter. Possible values must be between 1 and 50 characters in length and may only contain letters and numbers. Changing this forces a new resource to be created.

* `cdn_frontdoor_edge_action_id` - (Required) The ID of the Front Door Edge Action this execution filter belongs to. Changing this forces a new resource to be created.

* `cdn_frontdoor_edge_action_version_id` - (Required) The ID of the Front Door Edge Action Version which should run when a request matches this execution filter. The version must belong to the Edge Action specified in `cdn_frontdoor_edge_action_id`.

* `header_name` - (Required) The name of the request header which this execution filter matches.

* `header_value` - (Required) The value of the request header which this execution filter matches.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Edge Action Execution Filter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Edge Action Execution Filter.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Edge Action Execution Filter.
* `read` - (Defaults to 5 minutes) Used when retrieving the Edge Action Execution Filter.
* `update` - (Defaults to 30 minutes) Used when updating the Edge Action Execution Filter.
* `delete` - (Defaults to 30 minutes) Used when deleting the Edge Action Execution Filter.

## Import

Edge Action Execution Filters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_edge_action_execution_filter.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1/executionFilters/filter1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2026-10-01
