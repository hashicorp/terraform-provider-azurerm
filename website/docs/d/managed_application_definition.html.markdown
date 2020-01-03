---
subcategory: "ManagedApplication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application_definition"
sidebar_current: "docs-azurerm-datasource-managed-application-definition"
description: |-
Gets information about an existing Managed Application Definition
---

# Data Source: azurerm_managed_application_definition

Uses this data source to access information about an existing Managed Application Definition.

## Managed Application Definition Usage

```hcl
data "azurerm_managed_application_definition" "example" {
  resource_group_name = "acctestRG"
  name                = "acctestmanagedapplicationdefinition"
}

output "managed_application_definition_id" {
  value = "${data.azurerm_managed_application_definition.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Application Definition.

* `resource_group_name` - (Required) The Name of the Resource Group where the Managed Application Definition exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the Managed Application Definition exists.

* `authorization` - One or more `authorization` block defined below.

* `display_name` - The managed application definition display name.

* `lock_level` - The managed application lock level. Valid values include `ReadOnly`, `None`. Changing this forces a new resource to be created.

* `create_ui_definition` - The createUiDefinition json for the backing template with Microsoft.Solutions/applications resource.

* `description` - The managed application definition description.

* `main_template` - The inline main template json which has resources to be provisioned.

* `package_file_uri` - The managed application definition package file Uri.

* `tags` - A mapping of tags to assign to the resource.

---

An `authorization` block supports the following:

* `service_principal_id` - The provider's principal identifier. This is the identity that the provider will use to call ARM to manage the managed application resources.

* `role_definition_id` - The provider's role definition identifier. This role will define all the permissions that the provider must have on the managed application's container resource group. This role definition cannot have permission to delete the resource group.
