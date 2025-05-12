---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_gateway"
description: |-
  Manages a Spring Cloud Gateway.
---

# azurerm_spring_cloud_gateway

-> **Note:** This resource is applicable only for Spring Cloud Service with enterprise tier.

Manages a Spring Cloud Gateway.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_gateway` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_gateway" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id

  https_only                    = false
  public_network_access_enabled = true
  instance_count                = 2

  api_metadata {
    description       = "example description"
    documentation_url = "https://www.example.com/docs"
    server_url        = "https://wwww.example.com"
    title             = "example title"
    version           = "1.0"
  }

  cors {
    credentials_allowed = false
    allowed_headers     = ["*"]
    allowed_methods     = ["PUT"]
    allowed_origins     = ["example.com"]
    exposed_headers     = ["x-example-header"]
    max_age_seconds     = 86400
  }

  quota {
    cpu    = "1"
    memory = "2Gi"
  }

  sso {
    client_id     = "example id"
    client_secret = "example secret"
    issuer_uri    = "https://www.test.com/issueToken"
    scope         = ["read"]
  }

  local_response_cache_per_instance {
    size         = "100MB"
    time_to_live = "30s"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Gateway. Changing this forces a new Spring Cloud Gateway to be created. The only possible value is `default`.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new Spring Cloud Gateway to be created.

---

* `api_metadata` - (Optional) A `api_metadata` block as defined below.

* `application_performance_monitoring_ids` - (Optional) Specifies a list of Spring Cloud Application Performance Monitoring IDs.

* `application_performance_monitoring_types` - (Optional) Specifies a list of application performance monitoring types used in the Spring Cloud Gateway. The allowed values are `AppDynamics`, `ApplicationInsights`, `Dynatrace`, `ElasticAPM` and `NewRelic`.

* `client_authorization` - (Optional) A `client_authorization` block as defined below.

* `cors` - (Optional) A `cors` block as defined below.

* `environment_variables` - (Optional) Specifies the environment variables of the Spring Cloud Gateway as a map of key-value pairs.

* `https_only` - (Optional) is only https is allowed?

* `instance_count` - (Optional) Specifies the required instance count of the Spring Cloud Gateway. Possible Values are between `1` and `500`. Defaults to `1` if not specified.

* `public_network_access_enabled` - (Optional) Indicates whether the Spring Cloud Gateway exposes endpoint.

* `quota` - (Optional) A `quota` block as defined below.

* `local_response_cache_per_instance` - (Optional) A `local_response_cache_per_instance` block as defined below. Only one of `local_response_cache_per_instance` or `local_response_cache_per_route` can be specified.

* `local_response_cache_per_route` - (Optional) A `local_response_cache_per_route` block as defined below. Only one of `local_response_cache_per_instance` or `local_response_cache_per_route` can be specified.

* `sensitive_environment_variables` - (Optional) Specifies the sensitive environment variables of the Spring Cloud Gateway as a map of key-value pairs.

* `sso` - (Optional) A `sso` block as defined below.

---

A `api_metadata` block supports the following:

* `description` - (Optional) Detailed description of the APIs available on the Gateway instance.

* `documentation_url` - (Optional) Location of additional documentation for the APIs available on the Gateway instance.

* `server_url` - (Optional) Base URL that API consumers will use to access APIs on the Gateway instance.

* `title` - (Optional) Specifies the title describing the context of the APIs available on the Gateway instance.

* `version` - (Optional) Specifies the version of APIs available on this Gateway instance.

---

A `client_authorization` block supports the following:

* `certificate_ids` - (Optional) Specifies the Spring Cloud Certificate IDs of the Spring Cloud Gateway.

* `verification_enabled` - (Optional) Specifies whether the client certificate verification is enabled.

---

A `cors` block supports the following:

* `credentials_allowed` - (Optional) is user credentials are supported on cross-site requests?

* `allowed_headers` - (Optional) Allowed headers in cross-site requests. The special value `*` allows actual requests to send any header.

* `allowed_methods` - (Optional) Allowed HTTP methods on cross-site requests. The special value `*` allows all methods. If not set, `GET` and `HEAD` are allowed by default. Possible values are `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS` and `PUT`.

* `allowed_origins` - (Optional) Allowed origins to make cross-site requests. The special value `*` allows all domains.

* `allowed_origin_patterns` - (Optional) Allowed origin patterns to make cross-site requests.

* `exposed_headers` - (Optional) HTTP response headers to expose for cross-site requests.

* `max_age_seconds` - (Optional) How long, in seconds, the response from a pre-flight request can be cached by clients.

---

A `local_response_cache_per_route` block supports the following:

* `size` - (Optional) Specifies the maximum size of cache (10MB, 900KB, 1GB...) to determine if the cache needs to evict some entries.

* `time_to_live` - (Optional) Specifies the time before a cached entry is expired (300s, 5m, 1h...).

---

A `local_response_cache_per_instance` block supports the following:

* `size` - (Optional) Specifies the maximum size of cache (10MB, 900KB, 1GB...) to determine if the cache needs to evict some entries.

* `time_to_live` - (Optional) Specifies the time before a cached entry is expired (300s, 5m, 1h...).

---

The `quota` block supports the following:

* `cpu` - (Optional) Specifies the required cpu of the Spring Cloud Deployment. Possible Values are `500m`, `1`, `2`, `3` and `4`. Defaults to `1` if not specified.

-> **Note:** `cpu` supports `500m` and `1` for Basic tier, `500m`, `1`, `2`, `3` and `4` for Standard tier.

* `memory` - (Optional) Specifies the required memory size of the Spring Cloud Deployment. Possible Values are `512Mi`, `1Gi`, `2Gi`, `3Gi`, `4Gi`, `5Gi`, `6Gi`, `7Gi`, and `8Gi`. Defaults to `2Gi` if not specified.

-> **Note:** `memory` supports `512Mi`, `1Gi` and `2Gi` for Basic tier, `512Mi`, `1Gi`, `2Gi`, `3Gi`, `4Gi`, `5Gi`, `6Gi`, `7Gi`, and `8Gi` for Standard tier.

---

A `sso` block supports the following:

* `client_id` - (Optional) The public identifier for the application.

* `client_secret` - (Optional) The secret known only to the application and the authorization server.

* `issuer_uri` - (Optional) The URI of Issuer Identifier.

* `scope` - (Optional) It defines the specific actions applications can be allowed to do on a user's behalf.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Gateway.

* `url` - URL of the Spring Cloud Gateway, exposed when 'public_network_access_enabled' is true.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Gateway.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Gateway.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Gateway.

## Import

Spring Cloud Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_gateway.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1
```
