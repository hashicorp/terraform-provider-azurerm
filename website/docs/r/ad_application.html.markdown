---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ad_application"
sidebar_current: "docs-azurerm-resource-authorization-ad-application"
description: |-
  Manage an Azure Active Directory Application.

---

# azurerm_ad_application

Create a new application in Azure Active Directory. If your account is not an administrator in Active Directory an administrator must enable users to register applications within the User Settings. In addition, if you are using a Service Principal then it must have the permissions `Read and write all applications` and `Sign in and read user profile` under the `Windows Azure Active Directory` API.

## Example Usage

```hcl
resource "azurerm_ad_application" "example" {
  display_name = "example"
}
```

## Example Usage

```hcl
resource "azurerm_ad_application" "example" {
  display_name = "example"
  homepage = "http://homepage"
  identifier_uris = ["http://uri"]
  reply_urls = ["http://replyurl"]
  available_to_other_tenants = false
  oauth2_allow_implicit_flow = true
}
```

## Example Usage with Key Credential Certificate Authentication

```hcl
resource "tls_private_key" "example" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_self_signed_cert" "example" {
  key_algorithm   = "${tls_private_key.example.algorithm}"
  private_key_pem = "${tls_private_key.example.private_key_pem}"

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }

  validity_period_hours = 12

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
    "cert_signing",
  ]
}

resource "azurerm_ad_application" "example" {
  display_name = "example"

  key_credential {
    key_id = "32efa455-6b0e-489d-b3a7-d12675b767a5"
    type = "AsymmetricX509Cert"
    usage = "Verify"
    value = "${replace(tls_self_signed_cert.example.cert_pem, "/(-{5}.+?-{5})|(\\n)/", "")}"
  }
}
```

## Example Usage with Password Credential Authentication

```hcl
resource "azurerm_ad_application" "test" {
  display_name = "example"

  password_credential {
    key_id = "32efa455-6b0e-489d-b3a7-d12675b767a5"
    value = "example"
    start_date = "2018-03-01T00:00:00+00:00"
    end_date = "2018-03-02T00:00:00+00:00"
  }
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name for the application.

* `homepage` - (optional) The URL to the application's home page.

* `identifier_uris` - (Optional) User-defined URI(s) that uniquely identify a Web application within its Azure AD tenant, or within a verified custom domain if the application is multi-tenant.`

* `reply_urls` - (Optional) Specifies the URLs that user tokens are sent to for sign in, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to.

* `available_to_other_tenants` - (Optional) True if the application is shared with other tenants; otherwise, false.

* `oauth2_allow_implicit_flow` - (Optional) Specifies whether this web application can request OAuth2.0 implicit flow tokens.

* `key_credential` - (Optional) A list of Key Credential blocks as referenced below.

* `password_credential` - (Optional) A list of Password Credential blocks as referenced below.

`key_credential` supports the following:

* `key_id` - (Required) The unique identifier (GUID) for the key.

* `type` - (Required) The type of key credential. Possible values are `AsymmetricX509Cert` and `Symmetric`.

* `usage` - (Required) A string that describes the purpose for which the key can be used. Possible values are `Sign` and `Verify`.

* `start_date` - (Optional) The date and time at which the credential becomes valid.

* `end_date` - (Optional) The date and time at which the credential expires.

* `value` - (Required) The certificate value of the credential.

`password_credential` supports the following:

* `key_id` - (Required) The unique identifier (GUID) for the key.

* `start_date` - (Optional) The date and time at which the credential becomes valid.

* `end_date` - (Required) The date and time at which the credential expires.

* `value` - (Required) The secret value of the credential.

## Attributes Reference

The following attributes are exported:

* `app_id` - The Application ID.

* `object_id` - The Application Object ID.

## Import

Azure Active Directory Applications can be imported using the `object id`, e.g.

```shell
terraform import azurerm_ad_application.test 00000000-0000-0000-0000-000000000000
```
