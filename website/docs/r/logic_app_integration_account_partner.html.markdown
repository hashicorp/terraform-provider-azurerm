---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_partner"
description: |-
  Manages a Logic App Integration Account Partner.
---

# azurerm_logic_app_integration_account_partner

Manages a Logic App Integration Account Partner.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "example" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_logic_app_integration_account_partner" "example" {
  name                     = "example-iap"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.example.name

  business_identity {
    qualifier = "ZZ"
    value     = "AA"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Partner. Changing this forces a new Logic App Integration Account Partner to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Partner should exist. Changing this forces a new Logic App Integration Account Partner to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new Logic App Integration Account Partner to be created.

* `business_identity` - (Required) A `business_identity` block as documented below.

* `metadata` - (Optional) A JSON mapping of any Metadata for this Logic App Integration Account Partner.

---

A `business_identity` block exports the following:

* `qualifier` - (Required) The authenticating body that provides unique business identities to organizations.

* `value` - (Required) The value that identifies the documents that your logic apps receive.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Partner.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Partner.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Partner.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Partner.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Partner.

## Import

Logic App Integration Account Partners can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_partner.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/partners/partner1
```
