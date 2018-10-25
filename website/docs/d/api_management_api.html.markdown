---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api"
sidebar_current: "docs-azurerm-datasource-azurerm-api-management-api-x"
description: |-
  Gets information about an existing API Management API.
---

# Data Source: azurerm_api_management_api

Use this data source to access information about an existing API Management API.

## Example Usage

```hcl
data "azurerm_api_management_api" "test" {
  name                = "search-api"
  service_name        = "search-api-management"
  resource_group_name = "search-service"
}

output "api_management_api_id" {
  value = "${data.azurerm_api_management_api.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the API Management API.

* `service_name` - (Required) The name of the API Management service which the API Management API belongs to.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management API exists.

## Attributes Reference

* `id` - The ID of the API Management API.

* `path` - Relative URL uniquely identifying this API and all of its resource paths within the API Management service instance. It is appended to the API endpoint base URL specified during the service instance creation to form a public URL for this API. The path cannot start or end with `/`.

* `service_url` - Absolute URL of the backend service implementing this API.

* `description` - Description of the API. May include HTML formatting tags.

* `protocols` - Describes on which protocols the operations in this API can be invoked.

* `subscription_key_parameter_names` - Describes the names of the header and query parameter names used to send in the subscription key. The `subscription_key_parameter_names` block is defined below.

* `soap_pass_through` - Make API Management expose a SOAP front end, instead of a HTTP front end.

* `revision` - Describes the Revision of the Api.

* `version` - Indicates the Version identifier of the API if the API is versioned.

* `version_set_id` - A resource identifier for the related ApiVersionSet.

* `is_current` - Indicates if the API revision is current api revision.

* `is_online` - Indicates if the API revision is accessible via the gateway.

---

A `subscription_key_parameter_names` block exports the following:

* `header` - Subscription key header name.

* `query` - Subscription key query string parameter name.
