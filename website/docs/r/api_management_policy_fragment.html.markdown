---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_policy_fragment"
description: |-
  Manages an Api Management Policy Fragment.
---

# azurerm_api_management_policy_fragment

Manages an Api Management Policy Fragment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_policy_fragment" "example" {
  api_management_id = azurerm_api_management.example.id
  name              = "example-policy-fragment"
  format            = "xml"
  value             = file("policy-fragment-1.xml")
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The id of the API Management Service. Changing this forces a new Api Management Policy Fragment to be created.

* `name` - (Required) The name which should be used for this Api Management Policy Fragment. Changing this forces a new Api Management Policy Fragment to be created.

* `value` - (Required) The value of the Policy Fragment.

~> **Note:** Be aware of the two format possibilities. If the `value` is not applied and continues to cause a diff the format could be wrong.

* `format` - (Optional) The format of the Policy Fragment. Possible values are `xml` or `rawxml`. Default is `xml`.

~> **Note:** The `value` property will be updated to reflect the corresponding format when `format` is updated.

---

* `description` - (Optional) The description for the Policy Fragment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Api Management Policy Fragment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Api Management Policy Fragment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Api Management Policy Fragment.
* `update` - (Defaults to 30 minutes) Used when updating the Api Management Policy Fragment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Api Management Policy Fragment.

## Import

Api Management Policy Fragments can be imported using the `resource id`, e.g.

~> **Note:** Due to the behaviour of the API, Api Management Policy Fragments can only be imported as `xml`, but can be updated to the desired format after importing.

```shell
terraform import azurerm_api_management_policy_fragment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/policyFragments/policyFragment1
