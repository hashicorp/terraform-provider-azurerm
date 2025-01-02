---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_sas"
description: |-
  Gets a Shared Access Signature (SAS Token) for an existing Storage Account.

---

# Data Source: azurerm_storage_account_sas

Use this data source to obtain a Shared Access Signature (SAS Token) for an existing Storage Account.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of an Azure Storage Account.

Note that this is an [Account SAS](https://docs.microsoft.com/rest/api/storageservices/constructing-an-account-sas)
and *not* a [Service SAS](https://docs.microsoft.com/rest/api/storageservices/constructing-a-service-sas).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroupName"
  location = "West Europe"
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

data "azurerm_storage_account_sas" "example" {
  connection_string = azurerm_storage_account.example.primary_connection_string
  https_only        = true
  signed_version    = "2022-11-02"

  resource_types {
    service   = true
    container = false
    object    = false
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2018-03-21T00:00:00Z"
  expiry = "2020-03-21T00:00:00Z"

  permissions {
    read    = true
    write   = true
    delete  = false
    list    = false
    add     = true
    create  = true
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

output "sas_url_query_string" {
  value = data.azurerm_storage_account_sas.example.sas
}
```

## Argument Reference

* `connection_string` - The connection string for the storage account to which this SAS applies. Typically directly from the `primary_connection_string` attribute of a terraform created `azurerm_storage_account` resource.
* `https_only` - (Optional) Only permit `https` access. If `false`, both `http` and `https` are permitted. Defaults to `true`.
* `ip_addresses` - (Optional) IP address, or a range of IP addresses, from which to accept requests. When specifying a range, note that the range is inclusive.  
* `signed_version` - (Optional) Specifies the signed storage service version to use to authorize requests made with this account SAS. Defaults to `2022-11-02`.
* `resource_types` - A `resource_types` block as defined below.
* `services` - A `services` block as defined below.
* `start` - The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.
* `expiry` - The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.

-> **NOTE:** The [ISO-8601 Time offset from UTC](https://en.wikipedia.org/wiki/ISO_8601#Time_offsets_from_UTC) is currently not supported by the service, which will result into 409 error.

* `permissions` - A `permissions` block as defined below.

---

`resource_types` is a set of `true`/`false` flags which define the storage account resource types that are granted
access by this SAS. This can be thought of as the scope over which the permissions apply. A `service` will have
larger scope (affecting all sub-resources) than `object`.

A `resource_types` block contains:

* `service` - Should permission be granted to the entire service?
* `container` - Should permission be granted to the container?
* `object` - Should permission be granted only to a specific object?

---

`services` is a set of `true`/`false` flags which define the storage account services that are granted access by this SAS.

A `services` block contains:

* `blob` - Should permission be granted to `blob` services within this storage account?
* `queue` - Should permission be granted to `queue` services within this storage account?
* `table` - Should permission be granted to `table` services within this storage account?
* `file` - Should permission be granted to `file` services within this storage account?

---

A `permissions` block contains:

* `read` - Should Read permissions be enabled for this SAS?
* `write` - Should Write permissions be enabled for this SAS?
* `delete` - Should Delete permissions be enabled for this SAS?
* `list` - Should List permissions be enabled for this SAS?
* `add` - Should Add permissions be enabled for this SAS?
* `create` - Should Create permissions be enabled for this SAS?
* `update` - Should Update permissions be enabled for this SAS?
* `process` - Should Process permissions be enabled for this SAS?
* `tag` - Should Get / Set Index Tags permissions be enabled for this SAS?
* `filter` - Should Filter by Index Tags permissions be enabled for this SAS?

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/rest/api/storageservices/constructing-an-account-sas)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Account Shared Access Signature (SAS).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SAS Token.
