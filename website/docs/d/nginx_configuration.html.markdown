---
subcategory: "NGINX"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_configuration"
description: |-
  Gets information about an existing Nginx Configuration.
---

# Data Source: azurerm_nginx_configuration

Use this data source to access information about an existing Nginx Configuration.

## Example Usage

```hcl
data "azurerm_nginx_configuration" "example" {
  nginx_deployment_id = azurerm_nginx_deployment.example.id
}

output "id" {
  value = data.azurerm_nginx_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `nginx_deployment_id` - (Required) The ID of the Nginx Deployment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Nginx Configuration.

* `config_file` - A `config_file` block as defined below.

* `package_data` - The package data for this configuration.

* `root_file` - The root file path of this Nginx Configuration.

---

A `config_file` block exports the following:

* `content` - The base-64 encoded contents of this configuration file.

* `virtual_path` - The path of this configuration file.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Configuration.
