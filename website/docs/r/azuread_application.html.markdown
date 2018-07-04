---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_application"
sidebar_current: "docs-azurerm-resource-active-directory-application"
description: |-
  Manage an Azure Active Directory Application.

---

# azurerm_azuread_application

Create a new application in Azure Active Directory. If your account is not an administrator in Active Directory an administrator must enable users to register applications within the User Settings. In addition, if you are using a Service Principal then it must have the following permissions under the `Windows Azure Active Directory` API:

* Read and write all applications
* Sign in and read user profile

## Example Usage

```hcl
resource "azurerm_azuread_application" "example" {
  name = "example"
}
```

## Example Usage

```hcl
resource "azurerm_azuread_application" "example" {
  name = "example"
  homepage = "http://homepage"
  identifier_uris = ["http://uri"]
  reply_urls = ["http://replyurl"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name for the application.

* `homepage` - (optional) The URL to the application's home page.

* `identifier_uris` - (Optional) User-defined URI(s) that uniquely identify a Web application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant.`

* `reply_urls` - (Optional) Specifies the URLs that user tokens are sent to for sign in, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to.

* `available_to_other_tenants` - (Optional) True if the application is shared with other tenants; otherwise, false.

* `oauth2_allow_implicit_flow` - (Optional) Specifies whether this web application can request OAuth2.0 implicit flow tokens.

## Attributes Reference

The following attributes are exported:

* `application_id` - The Application ID.

## Import

Azure Active Directory Applications can be imported using the `object id`, e.g.

```shell
terraform import azurerm_azuread_application.test 00000000-0000-0000-0000-000000000000
```
