---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_application"
sidebar_current: "docs-azurerm-resource-azuread-application"
description: |-
  Manages an Application within Azure Active Directory.

---

# azurerm_azuread_application

Manages an Application within Azure Active Directory.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The display name for the application.

* `homepage` - (optional) The URL to the application's home page. If no homepage is specified this defaults to `http://{name}`.

* `identifier_uris` - (Optional) A list of user-defined URI(s) that uniquely identify a Web application within it's Azure AD tenant, or within a verified custom domain if the application is multi-tenant.

* `reply_urls` - (Optional) A list of URLs that user tokens are sent to for sign in, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to.

* `available_to_other_tenants` - (Optional) Is this Azure AD Application available to other tenants? Defaults to `false`.

* `oauth2_allow_implicit_flow` - (Optional) Does this Azure AD Application allow OAuth2.0 implicit flow tokens? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `application_id` - The Application ID.

## Import

Azure Active Directory Applications can be imported using the `object id`, e.g.

```shell
terraform import azurerm_azuread_application.test 00000000-0000-0000-0000-000000000000
```
