---
subcategory: "Managed Application"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application_definition"
description: |-
  Manages a Managed Application Definition.
---

# azurerm_managed_application_definition

Manages a Managed Application Definition.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_application_definition" "example" {
  name                = "example-managedapplicationdefinition"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  lock_level          = "ReadOnly"
  package_file_uri    = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name        = "TestManagedApplicationDefinition"
  description         = "Test Managed Application Definition"

  authorization {
    service_principal_id = data.azurerm_client_config.current.object_id
    role_definition_id   = "a094b430-dad3-424d-ae58-13f72fd72591"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Application Definition. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Application Definition should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `authorization` - (Required) One or more `authorization` block defined below.

* `lock_level` - (Required) Specifies the managed application lock level. Valid values include `CanNotDelete`, `None`, `ReadOnly`. Changing this forces a new resource to be created.

* `create_ui_definition` - (Optional) Specifies the createUiDefinition json for the backing template with Microsoft.Solutions/applications resource.

* `display_name` - (Optional) Specifies the managed application definition display name.

* `description` - (Optional) Specifies the managed application definition description.

* `package_enabled` - (Optional) Is enabled for this package? Defaults to `true`.

* `main_template` - (Optional) Specifies the inline main template json which has resources to be provisioned.

* `package_file_uri` - (Optional) Specifies the managed application definition package file Uri.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `authorization` block supports the following:

* `role_definition_id` - (Required) Specifies the provider's role definition identifier. This role will define all the permissions that the provider must have on the managed application's container resource group. This role definition cannot have permission to delete the resource group.

* `service_principal_id` - (Required) Specifies the provider's principal identifier. This is the identity that the provider will use to call ARM to manage the managed application resources.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Application Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Application Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Application Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Application Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Application Definition.

## Import

Managed Application Definition can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_managed_application_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Solutions/applicationDefinitions/appDefinition1
```
