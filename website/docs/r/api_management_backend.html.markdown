---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_backend"
description: |-
  Manages a backend within an API Management Service.
---

# azurerm_api_management_backend

Manages a backend within an API Management Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_backend" "example" {
  name                = "example-backend"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  protocol            = "http"
  url                 = "https://backend.com/api"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management backend. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where this backend should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `protocol` - (Required) The protocol used by the backend host. Possible values are `http` or `soap`.

* `url` - (Required) The backend host URL should be specified in the format `"https://backend.com/api"`, avoiding trailing slashes (/) to minimize misconfiguration risks. Azure API Management instance will append the backend resource name to this URL. This URL typically serves as the `base-url` in the [`set-backend-service`](https://learn.microsoft.com/azure/api-management/set-backend-service-policy) policy, enabling seamless transitions from frontend to backend.

---

* `credentials` - (Optional) A `credentials` block as documented below.

* `description` - (Optional) The description of the backend.

* `proxy` - (Optional) A `proxy` block as documented below.

* `resource_id` - (Optional) The management URI of the backend host in an external system. This URI can be the ARM Resource ID of Logic Apps, Function Apps or API Apps, or the management endpoint of a Service Fabric cluster.

* `service_fabric_cluster` - (Optional) A `service_fabric_cluster` block as documented below.

* `title` - (Optional) The title of the backend.

* `tls` - (Optional) A `tls` block as documented below.

---

A `credentials` block supports the following:

* `authorization` - (Optional) An `authorization` block as defined below.

* `certificate` - (Optional) A list of client certificate thumbprints to present to the backend host. The certificates must exist within the API Management Service.

* `header` - (Optional) A mapping of header parameters to pass to the backend host. The keys are the header names and the values are a comma separated string of header values. This is converted to a list before being passed to the API.

* `query` - (Optional) A mapping of query parameters to pass to the backend host. The keys are the query names and the values are a comma separated string of query values. This is converted to a list before being passed to the API.

---

An `authorization` block supports the following:

* `parameter` - (Optional) The authentication Parameter value.

* `scheme` - (Optional) The authentication Scheme name.

---

A `proxy` block supports the following:

* `password` - (Optional) The password to connect to the proxy server.

* `url` - (Required) The URL of the proxy server.

* `username` - (Required) The username to connect to the proxy server.

---

A `service_fabric_cluster` block supports the following:

* `client_certificate_thumbprint` - (Optional) The client certificate thumbprint for the management endpoint.

* `client_certificate_id` - (Optional) The client certificate resource id for the management endpoint.

~> **Note:** At least one of `client_certificate_thumbprint`, and `client_certificate_id` must be set.

* `management_endpoints` - (Required) A list of cluster management endpoints.

* `max_partition_resolution_retries` - (Required) The maximum number of retries when attempting resolve the partition.

* `server_certificate_thumbprints` - (Optional) A list of thumbprints of the server certificates of the Service Fabric cluster.

* `server_x509_name` - (Optional) One or more `server_x509_name` blocks as documented below.

---

A `server_x509_name` block supports the following:

* `issuer_certificate_thumbprint` - (Required) The thumbprint for the issuer of the certificate.

* `name` - (Required) The common name of the certificate.

---

A `tls` block supports the following:

* `validate_certificate_chain` - (Optional) Flag indicating whether SSL certificate chain validation should be done when using self-signed certificates for the backend host.

* `validate_certificate_name` - (Optional) Flag indicating whether SSL certificate name validation should be done when using self-signed certificates for the backend host.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Backend.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Backend.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Backend.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Backend.

## Import

API Management backends can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_backend.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/backends/backend1
```
