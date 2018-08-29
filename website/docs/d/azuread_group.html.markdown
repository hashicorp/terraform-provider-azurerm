---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_group"
sidebar_current: "docs-azurerm-datasource-azuread-group"
description: |-
  Gets information about a group object within the Azure Active Directory.

---

# Data Source: azurerm_azuread_group

Gets information about a Group object within the Azure Active Directory.

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to both `Read directory data` within the `Windows Azure Active Directory` API.

## Example Usage (by Object ID)

```hcl
data "azurerm_azuread_group" "test_group" {
  object_id = "78722cfc-8946-11e8-95f1-2200ec79ad01"
}
```

## Example Usage (by Group Display Name)

```hcl
data "azurerm_azuread_group" "test_group" {
  name = "MyTestGroup"
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Optional) The ID of the Azure AD Group we want to lookup.

* `name` - (Optional) The ID of the Azure AD Group we want to loopup.

-> **NOTE:** At least one of `name` or `object_id` must be specified.

-> **WARNING:** `name` is not unique within Azure Active Directory. The data source will only return the first Group found.

## Attributes Reference

The following attributes are exported:

* `id` - The Object ID for the Azure AD Group.
* `object_id` - The Object ID for the Azure AD Group.
* `name` - The Display Name for the Azure AD Group.
