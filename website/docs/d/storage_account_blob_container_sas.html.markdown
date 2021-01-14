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
  location = "westus"
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

## Argument Reference

* `connection_string` - The connection string for the storage account to which this SAS applies. Typically directly from the `primary_connection_string` attribute of a terraform created `azurerm_storage_account` resource.

* `container_name` - Name of the container.

* `https_only` - (Optional) Only permit `https` access. If `false`, both `http` and `https` are permitted. Defaults to `true`.

* `ip_address` - (Optional) Single ipv4 address or range (connected with a dash) of ipv4 addresses.

* `start` - The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.

* `expiry` - The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.

* `permissions` - A `permissions` block as defined below.

* `cache_control` - (Optional) The `Cache-Control` response header that is sent when this SAS token is used.

* `content_disposition` - (Optional) The `Content-Disposition` response header that is sent when this SAS token is used.

* `content_encoding` - (Optional) The `Content-Encoding` response header that is sent when this SAS token is used.

* `content_language` - (Optional) The `Content-Language` response header that is sent when this SAS token is used.

* `content_type` - (Optional) The `Content-Type` response header that is sent when this SAS token is used.

---

A `permissions` block contains:


* `read` - Should Read permissions be enabled for this SAS?

* `add` - Should Add permissions be enabled for this SAS?

* `create` - Should Create permissions be enabled for this SAS?

* `write` - Should Write permissions be enabled for this SAS?

* `delete` - Should Delete permissions be enabled for this SAS?

* `list` - Should List permissions be enabled for this SAS?

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/create-service-sas)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Blob Container Shared Access Signature (SAS).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blob Container.
