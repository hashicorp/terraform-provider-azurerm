---
subcategory: "ManagedApplication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application_definition"
sidebar_current: "docs-azurerm-resource-managed-application-definition"
description: |-
  Manages a Managed Application Definition.
---

# azurerm_managed_application_definition

Manages a Managed Application Definition.

## Managed Application Definition Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_application_definition" "example" {
  name                 = "example-managedapplicationdefinition"
  location             = "${azurerm_resource_group.example.location}"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  lock_level           = "ReadOnly"
  package_file_uri     = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name         = "TestManagedApplicationDefinition"
  description          = "Test Managed Application Definition"

  authorization {
    service_principal_id = "a1ac7e8c-14d3-432d-8c9b-0f780f99ef1e"
    role_definition_id   = "a094b430-dad3-424d-ae58-13f72fd72591"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Application Definition. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the Managed Application Definition should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `authorization` - (Required) One or more `authorization` block defined below.

* `display_name` - (Optional) The managed application definition display name.

* `lock_level` - (Optional) The managed application lock level. Valid values include `ReadOnly`, `None`. Changing this forces a new resource to be created.

* `create_ui_definition` - (Optional) The createUiDefinition json for the backing template with Microsoft.Solutions/applications resource. It can be a JObject or well-formed JSON string.

* `description` - (Optional) The managed application definition description.

* `main_template` - (Optional) The inline main template json which has resources to be provisioned. It can be a JObject or well-formed JSON string.

* `package_file_uri` - (Optional) The managed application definition package file Uri. Use this element.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `authorization` block supports the following:

* `service_principal_id` - (Required) The provider's principal identifier. This is the identity that the provider will use to call ARM to manage the managed application resources.

* `role_definition_id` - (Required) The provider's role definition identifier. This role will define all the permissions that the provider must have on the managed application's container resource group. This role definition cannot have permission to delete the resource group.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Application Definition.

## Import

Managed Application Definition can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_managed_application_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Solutions/applicationDefinitions/appDefinition1
```
