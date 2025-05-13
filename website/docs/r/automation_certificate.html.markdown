---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_certificate"
description: |-
  Manages an Automation Certificate.
---

# azurerm_automation_certificate

Manages an Automation Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "account1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "Basic"
}

resource "azurerm_automation_certificate" "example" {
  name                    = "certificate1"
  resource_group_name     = azurerm_resource_group.example.name
  automation_account_name = azurerm_automation_account.example.name

  description = "This is an example certificate"
  base64      = filebase64("certificate.pfx")
  exportable  = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Certificate. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Certificate is created. Changing this forces a new resource to be created.

* `automation_account_name` - (Required) The name of the automation account in which the Certificate is created. Changing this forces a new resource to be created.

* `base64` - (Required) Base64 encoded value of the certificate. Changing this forces a new resource to be created.

* `description` - (Optional) The description of this Automation Certificate.

* `exportable` - (Optional) The is exportable flag of the certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Automation Certificate ID.

* `thumbprint` - The thumbprint for the certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automation Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Automation Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automation Certificate.

## Import

Automation Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_automation_certificate.certificate1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Automation/automationAccounts/account1/certificates/certificate1
```
