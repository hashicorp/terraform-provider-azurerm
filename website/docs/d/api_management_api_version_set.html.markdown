---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_version_set"
description: |-
  Gets information about an existing API Version Set within an existing API Management Service.
---

# Data Source: azurerm_api_management_api_version_set

Uses this data source to access information about an API Version Set within an API Management Service.

## Example Usage

```hcl
data "azurerm_api_management_api_version_set" "example" {
  resource_group_name = "example-resources"
  api_management_name = "example-api"
  name                = "example-api-version-set"
}

output "api_management_api_version_set_id" {
  value = data.azurerm_api_management_api_version_set.example.id
}
```

## Argument Reference

* `name` - The name of the API Version Set.

* `resource_group_name` - The name of the Resource Group in which the parent API Management Service exists.

* `api_management_name` - The name of the [API Management Service](api_management.html) where the API Version Set exists.

## Attributes Reference

* `id` - The ID of the API Version Set.

* `description` - The description of API Version Set.

* `display_name` - The display name of this API Version Set.

* `versioning_schema` - The value that determines where the API Version identifer will be located in a HTTP request.

* `version_header_name` - The name of the Header which should be read from Inbound Requests which defines the API Version.

* `version_query_name` - The name of the Query String which should be read from Inbound Requests which defines the API Version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Version Set.
