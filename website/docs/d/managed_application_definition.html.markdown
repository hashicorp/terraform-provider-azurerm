---
subcategory: "ManagedApplication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_application_definition"
description: |-
  Gets information about an existing Managed Application Definition
---

# Data Source: azurerm_managed_application_definition

Uses this data source to access information about an existing Managed Application Definition.

## Managed Application Definition Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_managed_application_definition" "example" {
  name                = "example-managedappdef"
  resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_managed_application_definition.existing.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Application Definition.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where this Managed Application Definition exists.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure location where the resource exists.

* `authorization` - One or more `authorization` block defined below.

* `create_ui_definition` - The createUiDefinition json for the backing template with Microsoft.Solutions/applications resource.

* `description` - The managed application definition description.

* `display_name` - The managed application definition display name.

* `enabled` - The value indicating whether the package is enabled or not.

* `lock_level` - The managed application lock level.

* `main_template` - The inline main template json which has resources to be provisioned.

* `package_file_uri` - The managed application definition package file Uri.

* `tags` - A mapping of tags to assign to the resource.

---

An `authorization` block supports the following:

* `role_definition_id` - The provider's role definition identifier.

* `service_principal_id` - The provider's principal identifier.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataBox.
