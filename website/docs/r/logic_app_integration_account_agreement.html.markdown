---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_logic_app_integration_account_agreement"
description: |-
  Manages a Logic App Integration Account Agreement.
---

# azurerm_logic_app_integration_account_agreement

Manages a Logic App Integration Account Agreement.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "example-ia"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Standard"
}

resource "azurerm_logic_app_integration_account_partner" "host" {
  name                     = "example-hostpartner"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }
}

resource "azurerm_logic_app_integration_account_partner" "guest" {
  name                     = "example-guestpartner"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.test.name

  business_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}

resource "azurerm_logic_app_integration_account_agreement" "test" {
  name                     = "example-agreement"
  resource_group_name      = azurerm_resource_group.example.name
  integration_account_name = azurerm_logic_app_integration_account.test.name
  agreement_type           = "AS2"
  host_partner_name        = azurerm_logic_app_integration_account_partner.host.name
  guest_partner_name       = azurerm_logic_app_integration_account_partner.guest.name
  content                  = file("testdata/integration_account_agreement_content_as2.json")

  host_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamNY"
  }

  guest_identity {
    qualifier = "AS2Identity"
    value     = "FabrikamDC"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Logic App Integration Account Agreement. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Logic App Integration Account Agreement should exist. Changing this forces a new resource to be created.

* `integration_account_name` - (Required) The name of the Logic App Integration Account. Changing this forces a new resource to be created.

* `agreement_type` - (Required) The type of the Logic App Integration Account Agreement. Possible values are `AS2`, `X12` and `Edifact`.

* `content` - (Required) The content of the Logic App Integration Account Agreement.

* `guest_identity` - (Required) A `guest_identity` block as documented below.

* `guest_partner_name` - (Required) The name of the guest Logic App Integration Account Partner.

* `host_identity` - (Required) A `host_identity` block as documented below.

* `host_partner_name` - (Required) The name of the host Logic App Integration Account Partner.

* `metadata` - (Optional) The metadata of the Logic App Integration Account Agreement.

---

A `guest_identity` block exports the following:

* `qualifier` - (Required) The authenticating body that provides unique guest identities to organizations.

* `value` - (Required) The value that identifies the documents that your logic apps receive.

---

A `host_identity` block exports the following:

* `qualifier` - (Required) The authenticating body that provides unique host identities to organizations.

* `value` - (Required) The value that identifies the documents that your logic apps receive.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Logic App Integration Account Agreement.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Logic App Integration Account Agreement.
* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Integration Account Agreement.
* `update` - (Defaults to 30 minutes) Used when updating the Logic App Integration Account Agreement.
* `delete` - (Defaults to 30 minutes) Used when deleting the Logic App Integration Account Agreement.

## Import

Logic App Integration Account Agreements can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_logic_app_integration_account_agreement.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Logic/integrationAccounts/account1/agreements/agreement1
```
