---
subcategory: "Managed Applications"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application"
description: |-
  Manages a Managed Application.
---

# azurerm_managed_application

Manages a Managed Application.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_application_definition" "example" {
  name                = "examplemanagedapplicationdefinition"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  lock_level          = "ReadOnly"
  package_file_uri    = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name        = "TestManagedAppDefinition"
  description         = "Test Managed App Definition"

  authorization {
    service_principal_id = data.azurerm_client_config.current.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.builtin.id)[length(split("/", data.azurerm_role_definition.builtin.id)) - 1]
  }
}

resource "azurerm_managed_application" "example" {
  name                        = "example-managedapplication"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infrastructureGroup"
  application_definition_id   = azurerm_managed_application_definition.example.id

  parameter_values = jsonencode({
    location = {
      value = azurerm_resource_group.example.location
    },
    storageAccountNamePrefix = {
      value = "storeNamePrefix"
    },
    storageAccountType = {
      value = "Standard_LRS"
    }
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Application. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Application should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `kind` - (Required) The kind of the managed application to deploy. Possible values are `MarketPlace` and `ServiceCatalog`. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Required) The name of the target resource group where all the resources deployed by the managed application will reside. Changing this forces a new resource to be created.

* `application_definition_id` - (Optional) The application definition ID to deploy.

* `parameter_values` - (Optional) The parameter values to pass to the Managed Application. This field is a JSON object that allows you to assign parameters to this Managed Application.

* `plan` - (Optional) One `plan` block as defined below. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `plan` block exports the following:

* `name` - (Required) Specifies the name of the plan from the marketplace. Changing this forces a new resource to be created.

* `product` - (Required) Specifies the product of the plan from the marketplace. Changing this forces a new resource to be created.

* `publisher` - (Required) Specifies the publisher of the plan. Changing this forces a new resource to be created.

* `version` - (Required) Specifies the version of the plan from the marketplace. Changing this forces a new resource to be created.

* `promotion_code` - (Optional) Specifies the promotion code to use with the plan. Changing this forces a new resource to be created.

~> **Note:** When `plan` is specified, legal terms must be accepted for this item on this subscription before creating the Managed Application. The `azurerm_marketplace_agreement` resource or AZ CLI tool can be used to do this.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Application.

* `outputs` - The name and value pairs that define the managed application outputs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Application.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Application.

## Import

Managed Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Solutions/applications/app1
```
