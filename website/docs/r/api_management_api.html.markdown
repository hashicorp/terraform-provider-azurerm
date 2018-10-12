---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api"
sidebar_current: "docs-azurerm-resource-api-management-api-x"
description: |-
  Create a API Management API.
---

# azurerm_api_management_api

Create a API Management API component.

## Example Usage (import from Open API spec)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-dev"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "api-mngmnt-dev"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"
  sku {
    name = "Developer"
  }
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_api_management_api" "test" {
  name                = "conferenceapi"
  service_name        = "${azurerm_api_management.test.name}"
  path = "/conference"
  import {
    content_format = "swagger-link-json"
    content_value  = "http://conferenceapi.azurewebsites.net/?format=json"
  }
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management API.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management API exists.

* `location` - (Required) The Azure location where the API Management API exists.

* `service_name` - (Required) The Name of the API Management Service where the API Management API exists

* `path` - (Required) Relative URL uniquely identifying this API and all of its resource paths within the API Management service instance. It is appended to the API endpoint base URL specified during the service instance creation to form a public URL for this API.

---

* `service_url` - (Optional) Absolute URL of the backend service implementing this API.

* `description` - (Optional) Description of the API. May include HTML formatting tags.

* `protocols` - (Optional) A list of protocols the operations in this API can be invoked. Supported values are `http` and `https`. Default is `https`.

* `revision` - (Optional) The revision number of the API.

* `revision_description` - (Optional) Description of current revision.

-> **Note:** Setting revision using this resource is supported, but revision are more of a deployment concept than infrastructure, so we recommend finding other means to manage revisions.

* `import` - (Optional) A `import` block as documented below.

* `oauth` - (Optional) A `oauth` block as documented below.

* `subscription_key` - (Optional) A `subscription_key` block as documented below.

* `soap_api_type` - (Optional) Type of Soap API. Possible values include: 'http' or 'soap'. `http` creates a SOAP to REST API. `soap` creates a SOAP pass-through API. Default behavior when not set is REST API to REST API.

---

`import` block supports the following:

* `content_format` - (Required) Format of the Content in which the API is getting imported. Possible values include: 'swagger-json', 'swagger-link-json', 'wadl-link-json', 'wadl-xml', 'wsdl', 'wsdl-link'.

* `content_value` - (Required) Content value when Importing an API. When a `*-link-*` `content_format` is used, the `content_value` must be a URL. If not, `content_value` is defined inline.

* `wsdl_selector` - (Optional) Criteria to limit import of WSDL to a subset of the document. Only applicable to content with format `wsdl` or `wsdl-link`. The `wsdl_selector` block is documented below.

---

`wsdl_selector` block supports the following:

* `service_name` - (Required) Name of service to import from WSDL.

* `endpoint_name` - (Required) Name of endpoint(port) to import from WSDL.

---

`oauth` block supports the following:

* `authorization_server_id` - (Required) OAuth authorization server identifier.

* `scope` - (Required) Operations scope.

---

`subscription_key` block supports the following:

* `header` - (Optional) Subscription key header name.

* `query` - (Optional) Subscription key query string parameter name.

-> **Note:** Set both `header` and `query` to support using subscription key in both, or one of `header` or `query` to support one of them.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the API Management API component.

## Import

Api Management API can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/apis/api1
```
