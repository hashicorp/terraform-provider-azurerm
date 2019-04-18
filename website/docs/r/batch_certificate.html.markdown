---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_certificate"
sidebar_current: "docs-azurerm-resource-batch-certificate"
description: |-
  Manages a certificate in an Azure Batch account.

---

# azurerm_batch_certificate

Manages a certificate in an Azure Batch account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testbatch"
  location = "westeurope"
}

resource "azurerm_storage_account" "test" {
  name                     = "teststorage"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testbatchaccount"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_certificate" "test" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${base64encode(file("certificate.pfx"))}"
  format               = "Pfx"
  password             = "terraform"
  thumbprint           = "42C107874FD0E4A9583292A2F1098E8FE4B2EDDA"
  thumbprint_algorithm = "SHA1"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) Specifies the name of the Batch account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch account. Changing this forces a new resource to be created.

* `certificate` - (Required) The base64-encoded contents of the certificate.

* `format` - (Required) The format of the certificate. Possible values are `Cer` or `Pfx`.

* `password` - (Optional) The password to access the certificate's private key. This must and can only be specified when `format` is `Pfx`.

* `thumbprint` - (Required) The thumbprint of the certificate. At this time the only supported value is 'SHA1'.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch certificate ID.

* `name` - The generated name of the certificate.

* `public_data` - The public key of the certificate.
