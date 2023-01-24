---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_certificate"
description: |-
  Gets information about a Container App Environment Certificate.
---

# Data Source: azurerm_container_app_environment_certificate.

Use this data source to access information about an existing Container App Environment Certificate.

## Example Usage

```hcl
data "azurerm_container_app_environment" "example" {
  name                = "example-environment"
  resource_group_name = "example-resources"
}

data "azurerm_container_app_environment_certificate" "example" {
  name                         = "mycertificate"
  container_app_environment_id = data.azurerm_container_app_environment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Certificate. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The ID of the Container App Environment to configure this Certificate on. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Certificate

* `expiration_date` - The expiration date for the Certificate.

* `issue_date` - The date of issue for the Certificate.

* `issuer` - The Certificate Issuer.

* `subject_name` - The Subject Name for the Certificate.

* `thumbprint` - The Thumbprint of the Certificate.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Certificate.
