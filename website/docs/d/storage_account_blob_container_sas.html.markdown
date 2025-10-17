---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_blob_container_sas"
description: |-
  Gets a Shared Access Signature (SAS Token) for an existing Storage Account Blob Container.

---

# Data Source: azurerm_storage_account_blob_container_sas

Use this data source to obtain a Shared Access Signature (SAS Token) for an existing Storage Account Blob Container.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of an Azure Storage Account Blob Container.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroupName"
  location = "West Europe"
}

resource "azurerm_storage_account" "storage" {
  name                     = "storageaccountname"
  resource_group_name      = azurerm_resource_group.rg.name
  location                 = azurerm_resource_group.rg.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "container" {
  name                  = "mycontainer"
  storage_account_name  = azurerm_storage_account.storage.name
  container_access_type = "private"
}

data "azurerm_storage_account_blob_container_sas" "example" {
  connection_string = azurerm_storage_account.storage.primary_connection_string
  container_name    = azurerm_storage_container.container.name
  https_only        = true

  ip_address = "168.1.5.65"

  start  = "2018-03-21"
  expiry = "2018-03-21"

  permissions {
    read   = true
    add    = true
    create = false
    write  = false
    delete = true
    list   = true
  }

  cache_control       = "max-age=5"
  content_disposition = "inline"
  content_encoding    = "deflate"
  content_language    = "en-US"
  content_type        = "application/json"
}

output "sas_url_query_string" {
  value = data.azurerm_storage_account_blob_container_sas.example.sas
}
```

## Arguments Reference

* `connection_string` - (Required) The connection string for the storage account to which this SAS applies. Typically directly from the `primary_connection_string` attribute of a terraform created `azurerm_storage_account` resource.

* `container_name` - (Required) Name of the container.

* `https_only` - (Optional) Only permit `https` access. If `false`, both `http` and `https` are permitted. Defaults to `true`.

* `ip_address` - (Optional) Single IPv4 address or range (connected with a dash) of IPv4 addresses.

* `start` - (Required) The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.

* `expiry` - (Required) The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.

-> **Note:** The [ISO-8601 Time offset from UTC](https://en.wikipedia.org/wiki/ISO_8601#Time_offsets_from_UTC) is currently not supported by the service, which will result into 409 error.

* `permissions` - (Required) A `permissions` block as defined below.

* `cache_control` - (Optional) The `Cache-Control` response header that is sent when this SAS token is used.

* `content_disposition` - (Optional) The `Content-Disposition` response header that is sent when this SAS token is used.

* `content_encoding` - (Optional) The `Content-Encoding` response header that is sent when this SAS token is used.

* `content_language` - (Optional) The `Content-Language` response header that is sent when this SAS token is used.

* `content_type` - (Optional) The `Content-Type` response header that is sent when this SAS token is used.

---

A `permissions` block contains:

* `read` - (Optional) Should Read permissions be enabled for this SAS?

* `add` - (Optional) Should Add permissions be enabled for this SAS?

* `create` - (Optional) Should Create permissions be enabled for this SAS?

* `write` - (Optional) Should Write permissions be enabled for this SAS?

* `delete` - (Optional) Should Delete permissions be enabled for this SAS?

* `delete_version` - (Optional) Should Delete version permissions be enabled for this SAS?

* `list` - (Optional) Should List permissions be enabled for this SAS?

* `tags` - (Optional) Should Tags permissions be enabled for this SAS?

* `find` - (Optional) Should Find permissions be enabled for this SAS?

* `move` - (Optional) Should Move permissions be enabled for this SAS?

* `execute` - (Optional) Should Execute permissions be enabled for this SAS?

* `ownership` - (Optional) Should Ownership permissions be enabled for this SAS?

* `permissions` - (Optional) Should Permissions permissions be enabled for this SAS?

* `set_immutability_policy` - (Optional) Should Set Immutability Policy permissions be enabled for this SAS?


Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/rest/api/storageservices/create-service-sas)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Blob Container Shared Access Signature (SAS). The delimiter character ('?') for the query string is the prefix of `sas`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blob Container.
