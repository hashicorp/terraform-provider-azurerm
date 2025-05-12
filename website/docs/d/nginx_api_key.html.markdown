---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_api_key"
description: |-
  Gets information about an existing NGINX Dataplane API Key.
---

# Data Source: azurerm_nginx_api_key

Use this data source to access information about an existing NGINX Dataplane API Key.

## Example Usage

```hcl
data "azurerm_nginx_api_key" "example" {
  name                = "existing"
  nginx_deployment_id = azurerm_nginx_deployment.example.id
}

output "id" {
  value = data.azurerm_nginx_api_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the NGINX Dataplane API Key.

- `nginx_deployment_id` - (Required) The ID of the NGINX Deployment that the API key is associated with.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `end_date_time` - The RFC3339 formatted time after which this Dataplane API Key is no longer valid.

- `hint` - The first three characters of the secret text to help identify it in use.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NGINX Dataplane API Key.
