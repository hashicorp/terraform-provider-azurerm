---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_custom_domain"
description: |-
  Manages a Container App Environment Custom Domain.
---

# azurerm_container_app_environment_custom_domain

Manages a Container App Environment Custom Domain Suffix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "my-environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_custom_domain" "example" {
  container_app_environment_id = azurerm_container_app_environment.example.id
  certificate_blob_base64      = filebase64("testacc.pfx")
  certificate_password         = "TestAcc"
  dns_suffix                   = "acceptancetest.contoso.com"
}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The ID of the Container Apps Managed Environment. Changing this forces a new resource to be created.

* `certificate_blob_base64` - (Required) The bundle of Private Key and Certificate for the Custom DNS Suffix as a base64 encoded PFX or PEM.

* `certificate_password` - (Required) The password for the Certificate bundle.

* `dns_suffix` - (Required) Custom DNS Suffix for the Container App Environment.
## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment.

## Import

A Container App Environment Custom Domain Suffix can be imported using the `resource id` of its parent container ontainer App Environment , e.g.

```shell
terraform import azurerm_container_app_environment_custom_domain.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myEnvironment"
```
