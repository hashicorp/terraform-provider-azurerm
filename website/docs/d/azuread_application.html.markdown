---
subcategory: "Azure Active Directory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_application"
sidebar_current: "docs-azurerm-datasource-azuread-application"
description: |-
  Gets information about an existing Application within Azure Active Directory.
---

# Data Source: azurerm_azuread_application

Use this data source to access information about an existing Application within Azure Active Directory.

~> **NOTE:** The Azure Active Directory resources have been split out into [a new AzureAD Provider](http://terraform.io/docs/providers/azuread/index.html) - as such the AzureAD resources within the AzureRM Provider are deprecated and will be removed in the next major version (2.0). Information on how to migrate from the existing resources to the new AzureAD Provider [can be found here](../guides/migrating-to-azuread.html).

-> **NOTE:** If you're authenticating using a Service Principal then it must have permissions to both `Read and write all applications` and `Sign in and read user profile` within the `Windows Azure Active Directory` API.

## Example Usage

```hcl
data "azurerm_azuread_application" "example" {
  name = "My First AzureAD Application"
}

output "azure_active_directory_object_id" {
  value = "${data.azurerm_azuread_application.example.id}"
}
```

## Argument Reference

* `object_id` - (Optional) Specifies the Object ID of the Application within Azure Active Directory.

* `name` - (Optional) Specifies the name of the Application within Azure Active Directory.

-> **NOTE:** Either an `object_id` or `name` must be specified.

## Attributes Reference

* `id` - the Object ID of the Azure Active Directory Application.

* `application_id` - the Application ID of the Azure Active Directory Application.

* `available_to_other_tenants` - Is this Azure AD Application available to other tenants?

* `identifier_uris` - A list of user-defined URI(s) that uniquely identify a Web application within it's Azure AD tenant, or within a verified custom domain if the application is multi-tenant.

* `oauth2_allow_implicit_flow` - Does this Azure AD Application allow OAuth2.0 implicit flow tokens?

* `object_id` - the Object ID of the Azure Active Directory Application.

* `reply_urls` - A list of URLs that user tokens are sent to for sign in, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to.
