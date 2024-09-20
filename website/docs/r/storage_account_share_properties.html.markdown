---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_share_properties"
description: |-
  Manages a Azure Storage Account Share Properties.
---

# azurerm_storage_account_share_properties

Manages an Azure Storage Accounts Share Properties.

## Disclaimers

~> **Note on Storage Accounts and Share Properties:** Terraform currently provides both a standalone [Share Properties resource](storage_account_share_properties.html), and allows for Share Properties to be defined in-line within the [Storage Account resource](storage_account.html). At this time you cannot use a Storage Account with in-line Share Properties in conjunction with any Share Properties resource. Doing so will cause a conflict of Share Properties configurations and will overwrite the in-line Share Properties.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_storage_account" "example" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account_share_properties" "example" {
  storage_account_id = azurerm_storage_account.example.id

  cors_rule {
    allowed_origins    = ["http://www.example.com"]
    exposed_headers    = ["x-tempo-*"]
    allowed_headers    = ["x-tempo-*"]
    allowed_methods    = ["GET", "PUT", "PATCH"]
    max_age_in_seconds = "500"
  }

  retention_policy {
    days = 300
  }

  smb {
    versions                        = ["SMB3.0"]
    authentication_types            = ["NTLMv2"]
    kerberos_ticket_encryption_type = ["AES-256"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) Specifies the resource id of the storage account.

* `cors_rule` - (Optional) A `cors_rule` block as defined below.

* `retention_policy` - (Optional) A `retention_policy` block as defined below.

* `smb` - (Optional) A `smb` block as defined below.

---

A `cors_rule` block supports the following:

* `allowed_headers` - (Required) A list of headers that are allowed to be a part of the cross-origin request.

* `allowed_methods` - (Required) A list of HTTP methods that are allowed to be executed by the origin. Valid options are `DELETE`, `GET`, `HEAD`, `MERGE`, `POST`, `OPTIONS`, `PUT` or `PATCH`.

* `allowed_origins` - (Required) A list of origin domains that will be allowed by CORS.

* `exposed_headers` - (Required) A list of response headers that are exposed to CORS clients.

* `max_age_in_seconds` - (Required) The number of seconds the client should cache a preflight response.

---

A `retention_policy` block supports the following:

* `days` - (Optional) Specifies the number of days that the `azurerm_storage_share` should be retained, between `1` and `365` days. Defaults to `7`.

---

A `smb` block supports the following:

* `versions` - (Optional) A set of SMB protocol versions. Possible values are `SMB2.1`, `SMB3.0`, and `SMB3.1.1`.

* `authentication_types` - (Optional) A set of SMB authentication methods. Possible values are `NTLMv2`, and `Kerberos`.

* `kerberos_ticket_encryption_type` - (Optional) A set of Kerberos ticket encryption. Possible values are `RC4-HMAC`, and `AES-256`.

* `channel_encryption_type` - (Optional) A set of SMB channel encryption. Possible values are `AES-128-CCM`, `AES-128-GCM`, and `AES-256-GCM`.

* `multichannel_enabled` - (Optional) Indicates whether multichannel is enabled. Defaults to `false`. This is only supported on Premium storage accounts.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Account Share Properties.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Storage Account Share Properties.
* `update` - (Defaults to 60 minutes) Used when updating the Storage Account Share Properties.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Share Properties.
* `delete` - (Defaults to 60 minutes) Used when deleting the Storage Account Share Properties.

## Import

Storage Account Share Properties can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account.share1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
