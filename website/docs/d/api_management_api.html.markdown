---
subcategory: "API Management"
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
data "azurerm_api_management_api" "example" {
  name                = "search-api"
  api_management_name = "search-api-management"
  resource_group_name = "search-service"
  revision            = "2"
}

output "api_management_api_id" {
  value = "${data.azurerm_api_management_api.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the API Management API.

* `api_management_name` - (Required) The name of the API Management Service in which the API Management API exists.

* `resource_group_name` - (Required) The Name of the Resource Group in which the API Management Service exists.

* `revision` - (Required) The Revision of the API Management API.

## Attributes Reference

* `id` - The ID of the API Management API.

* `description` - A description of the API Management API, which may include HTML formatting tags.

* `display_name` - The display name of the API.

* `is_current` - Is this the current API Revision?

* `is_online` - Is this API Revision online/accessible via the Gateway?

* `path` - The Path for this API Management API.

* `protocols` - A list of protocols the operations in this API can be invoked.

* `service_url` - Absolute URL of the backend service implementing this API.

* `soap_pass_through` - Should this API expose a SOAP frontend, rather than a HTTP frontend?

* `subscription_key_parameter_names` - A `subscription_key_parameter_names` block as documented below.

* `version` - The Version number of this API, if this API is versioned.

* `version_set_id` - The ID of the Version Set which this API is associated with.

---

A `subscription_key_parameter_names` block exports the following:

* `header` - The name of the HTTP Header which should be used for the Subscription Key.

* `query` - The name of the QueryString parameter which should be used for the Subscription Key.

---

A `wsdl_selector` block exports the following:

* `service_name` - The name of service to import from WSDL.

* `endpoint_name` - The name of endpoint (port) to import from WSDL.
