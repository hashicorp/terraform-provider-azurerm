---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_sas"
sidebar_current: "docs-azurerm-datasource-storage-account-sas"
description: |-
  Create a Shared Access Signature (SAS) for an Azure Storage Account.

---

# Data Source: azurerm_storage_account_sas

Use this data source to create a Shared Access Signature (SAS) for an Azure Storage Account.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of an Azure Storage Account.

Note that this is an [Account SAS](https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-an-account-sas)
and *not* a [Service SAS](https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-a-service-sas).

## Example Usage

```hcl
resource "azurerm_resource_group" "testrg" {
  name     = "resourceGroupName"
  location = "westus"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "storageaccountname"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "westus"
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags {
    environment = "staging"
  }
}

data "azurerm_storage_account_sas" "test" {
    connection_string = "${azurerm_storage_account.testsa.primary_connection_string}"
    https_only        = true
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
    start   = "2018-03-21"
    expiry  = "2020-03-21"
    permissions {
        read    = true
        write   = true
        delete  = false
        list    = false
        add     = true
        create  = true
        update  = false
        process = false
    }
}

output "sas_url_query_string" {
  value = "${data.azurerm_storage_account_sas.test.sas}"
}
```

## Argument Reference

* `connection_string` - (Required) The connection string for the storage account to which this SAS applies. Typically directly from the `primary_connection_string` attribute of a terraform created `azurerm_storage_account` resource.
* `https_only` - (Optional) Only permit `https` access. If `false`, both `http` and `https` are permitted. Defaults to `true`.
* `resouce_types` - (Required) A `resource_types` block as defined below. 
* `services` - (Required) A `services` block as defined below.
* `start` - (Required) The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.
* `expiry` - (Required) The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.
* `permissions` - (Required) A `permissions` block as defined below.

---

`resource_types` is a set of `true`/`false` flags which define the storage account resource types that are granted 
access by this SAS. This can be thought of as the scope over which the permissions apply. A `service` will have
larger scope (affecting all sub-resources) than `object`.

A `resource_types` block contains:

* `service` - (Required) Should permission be granted to the entire service?
* `container` - (Required) Should permission be granted to the container?
* `object` - (Required) Should permission be granted only to a specific object?

---

`services` is a set of `true`/`false` flags which define the storage account services that are granted access by this SAS.

A `services` block contains:

* `blob` - (Required) Should permission be granted to `blob` services within this storage account?
* `queue` - (Required) Should permission be granted to `queue` services within this storage account?
* `table` - (Required) Should permission be granted to `table` services within this storage account?
* `file` - (Required) Should permission be granted to `file` services within this storage account?

---

A `permissions` block contains:

* `read` - (Required) Should Read permissions be enabled for this SAS?
* `write` - (Required) Should Write permissions be enabled for this SAS?
* `delete` - (Required) Should Delete permissions be enabled for this SAS?
* `list` - (Required) Should List permissions be enabled for this SAS?
* `add` - (Required) Should Add permissions be enabled for this SAS?
* `create` - (Required) Should Create permissions be enabled for this SAS?
* `update` - (Required) Should Update permissions be enabled for this SAS?
* `process` - (Required) Should Process permissions be enabled for this SAS?

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-an-account-sas)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Account Shared Access Signature (SAS). 
