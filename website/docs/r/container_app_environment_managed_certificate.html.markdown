---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_managed_certificate"
description: |-
  Manages a Container App Environment Managed Certificate.
---

# azurerm_container_app_environment_managed_certificate

Manages a Container App Environment Managed Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "example-environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_managed_certificate" "example" {
  name                         = "example-managed-cert"
  container_app_environment_id = azurerm_container_app_environment.example.id
  subject_name                 = "example.com"
  domain_control_validation    = "HTTP"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Environment Managed Certificate. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Managed Certificate on. Changing this forces a new resource to be created.

* `subject_name` - (Required) The Subject Name of the Certificate. Must be a valid domain name. Changing this forces a new resource to be created.

---

* `domain_control_validation` - (Optional) The domain control validation type for the managed certificate. Possible values are `CNAME`, `HTTP` and `TXT`. Defaults to `HTTP`. Changing this forces a new resource to be created.

~> **Note:** The supported validation methods depend on the domain. Azure will validate domain ownership based on the specified method. `HTTP` validation requires an HTTP endpoint at the domain, `CNAME` validation requires DNS CNAME record configuration, and `TXT` validation requires DNS TXT record configuration.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Managed Certificate.

* `validation_token` - The validation token for the managed certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Managed Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Managed Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Managed Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Managed Certificate.

## Import

A Container App Environment Managed Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_managed_certificate.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myenv/managedCertificates/mycertificate"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
