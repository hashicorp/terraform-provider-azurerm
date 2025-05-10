---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_certificate"
description: |-
  Manages a Container App Environment Certificate.
---

# azurerm_container_app_environment_certificate

Manages a Container App Environment Certificate.

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
  name                       = "myEnvironment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_environment_certificate" "example" {
  name                         = "myfriendlyname"
  container_app_environment_id = azurerm_container_app_environment.example.id
  certificate_blob             = filebase64("path/to/certificate_file.pfx")
  certificate_password         = "$3cretSqu1rreL"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container Apps Environment Certificate. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The Container App Managed Environment ID to configure this Certificate on. Changing this forces a new resource to be created.

* `certificate_blob_base64` - (Required) The Certificate Private Key as a base64 encoded PFX or PEM. Changing this forces a new resource to be created.

* `certificate_password` - (Required) The password for the Certificate. Changing this forces a new resource to be created.

---

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Certificate

* `expiration_date` - The expiration date for the Certificate.

* `issue_date` - The date of issue for the Certificate.

* `issuer` - The Certificate Issuer.

* `subject_name` - The Subject Name for the Certificate.

* `thumbprint` - The Thumbprint of the Certificate.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Environment Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Environment Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Environment Certificate.

## Import

A Container App Environment Certificate can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app_environment_certificate.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/managedEnvironments/myenv/certificates/mycertificate"
```
