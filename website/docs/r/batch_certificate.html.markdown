---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_certificate"
description: |-
  Manages a certificate in an Azure Batch account.

---

# azurerm_batch_certificate

Manages a certificate in an Azure Batch account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "testbatch"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "teststorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "example" {
  name                 = "testbatchaccount"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.example.id

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_certificate" "example" {
  resource_group_name  = azurerm_resource_group.example.name
  account_name         = azurerm_batch_account.example.name
  certificate          = filebase64("certificate.pfx")
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

* `password` - (Optional) The password to access the certificate's private key. This can only be specified when `format` is `Pfx`.

* `thumbprint` - (Required) The thumbprint of the certificate. At this time the only supported value is 'SHA1'.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Batch Certificate.

* `name` - The generated name of the certificate.

* `public_data` - The public key of the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Batch Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Batch Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Batch Certificate.

## Import

Batch Certificates can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_batch_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Batch/batchAccounts/batch1/certificates/certificate1
```
