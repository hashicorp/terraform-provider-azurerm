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
  value = "${data.azurerm_storage_account_sas.sas}"
}
```

## Argument Reference

* `connection_string` - (Required) The connection string for the storage account to which this SAS applies. Typically directly from the `primary_connection_string` attribute of a terraform created `azurerm_storage_account` resource.
* `https_only` - (Optional) Only permit `https` access. If `false`, both `http` and `https` are permitted. Defaults to `true`.
* `resouce_types` - (Required) A set of `true`/`false` flags which define the storage account resource types that are granted access by this SAS. 
* `services` - (Required) A set of `true`/`false` flags which define the storage account services that are granted access by this SAS.
* `start` - (Required) The starting time and date of validity of this SAS. Must be a valid ISO-8601 format time/date string.
* `expiry` - (Required) The expiration time and date of this SAS. Must be a valid ISO-8601 format time/date string.
* `permissions` - (Required) A set of `true`/`false` flags which define the granted permissions.

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/en-us/rest/api/storageservices/constructing-an-account-sas)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Account Shared Access Signature (SAS). 
