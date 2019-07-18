---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_azuread_service_principal_password"
sidebar_current: "docs-azurerm-resource-azuread-service-principal-password"
description: |-
  Manages a Password associated with a Service Principal within Azure Active Directory.

---

# azurerm_azuread_service_principal_password

Manages a Password associated with a Service Principal within Azure Active Directory.

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

resource "azurerm_azuread_service_principal_password" "test" {
  service_principal_id = "${azurerm_azuread_service_principal.test.id}"
  value                = "VT=uSgbTanZhyz@%nL9Hpd+Tfay_MRV#"
  end_date             = "2020-01-01T01:02:03Z"
}
```

## Argument Reference

The following arguments are supported:

* `service_principal_id` - (Required) The ID of the Service Principal for which this password should be created. Changing this field forces a new resource to be created.

* `value` - (Required) The Password for this Service Principal.

* `end_date` - (Required) The End Date which the Password is valid until, formatted as a RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). Changing this field forces a new resource to be created.

* `key_id` - (Optional) A GUID used to uniquely identify this Key. If not specified a GUID will be created. Changing this field forces a new resource to be created.

* `start_date` - (Optional) The Start Date which the Password is valid from, formatted as a RFC3339 date string (e.g. `2018-01-01T01:02:03Z`). If this isn't specified, the current date is used.  Changing this field forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The Key ID for the Service Principal Password.

## Import

Service Principal Passwords can be imported using the `object id`, e.g.

```shell
terraform import azurerm_azuread_service_principal_password.test 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111
```

-> **NOTE:** This ID format is unique to Terraform and is composed of the Service Principal's Object ID and the Service Principal Password's Key ID in the format `{ServicePrincipalObjectId}/{ServicePrincipalPasswordKeyId}`.
