---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_service_principal"
sidebar_current: "docs-azurerm-resource-azuread-service-principal-x"
description: |-
  Manages a Service Principal associated with an Application within Azure Active Directory.

---

# azurerm_azuread_service_principal

Manages a Service Principal associated with an Application within Azure Active Directory.

~> **NOTE:** The Azure Active Directory resources have been split out into [a new AzureAD Provider](http://terraform.io/docs/providers/azuread/index.html) - as such the AzureAD resources within the AzureRM Provider are deprecated and will be removed in the next major version (2.0). Information on how to migrate from the existing resources to the new AzureAD Provider [can be found here](../guides/migrating-to-azuread.html).

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to both `Read and write all applications` and `Sign in and read user profile` within the `Windows Azure Active Directory` API.

## Example Usage

```hcl
resource "azurerm_azuread_application" "test" {
  name                       = "example"
  homepage                   = "http://homepage"
  identifier_uris            = ["http://uri"]
  reply_urls                 = ["http://replyurl"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}

resource "azurerm_azuread_service_principal" "test" {
  application_id = "${azurerm_azuread_application.test.application_id}"
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The ID of the Azure AD Application for which to create a Service Principal.

## Attributes Reference

The following attributes are exported:

* `id` - The Object ID for the Service Principal.

* `display_name` - The Display Name of the Azure Active Directory Application associated with this Service Principal.

## Import

Azure Active Directory Service Principals can be imported using the `object id`, e.g.

```shell
terraform import azurerm_azuread_service_principal.test 00000000-0000-0000-0000-000000000000
```
