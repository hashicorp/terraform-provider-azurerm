---
subcategory: "Azure Active Directory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_service_principal"
sidebar_current: "docs-azurerm-datasource-azuread-service-principal"
description: |-
  Gets information about an existing Service Principal associated with an Application within Azure Active Directory.

---

# Data Source: azurerm_azuread_service_principal

Gets information about an existing Service Principal associated with an Application within Azure Active Directory.

~> **NOTE:** The Azure Active Directory resources have been split out into [a new AzureAD Provider](http://terraform.io/docs/providers/azuread/index.html) - as such the AzureAD resources within the AzureRM Provider are deprecated and will be removed in the next major version (2.0). Information on how to migrate from the existing resources to the new AzureAD Provider [can be found here](../guides/migrating-to-azuread.html).

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to both `Read and write all applications` and `Sign in and read user profile` within the `Windows Azure Active Directory` API.

## Example Usage (by Application Display Name)

```hcl
data "azurerm_azuread_service_principal" "example" {
  display_name = "my-awesome-application"
}
```

## Example Usage (by Application ID)

```hcl
data "azurerm_azuread_service_principal" "example" {
  application_id = "00000000-0000-0000-0000-000000000000"
}
```

## Example Usage (by Object ID)

```hcl
data "azurerm_azuread_service_principal" "example" {
  object_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Optional) The ID of the Azure AD Application for which to create a Service Principal.

* `object_id` - (Optional) The ID of the Azure AD Service Principal.

* `display_name` - (Optional) The Display Name of the Azure AD Application associated with this Service Principal.

-> **NOTE:** At least one of `application_id`, `display_name` or `object_id` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The Object ID for the Service Principal.
