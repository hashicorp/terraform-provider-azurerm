---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product"
sidebar_current: "docs-azurerm-resource-api-management-product"
description: |-
  Manages an API Management Product.
---

# azurerm_api_management_product

Manages an API Management Product.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "test" {
  name                = "example-apim"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = "${azurerm_api_management.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  published             = true
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `approval_required` - (Optional) Do subscribers need to be approved prior to being able to use the Product?

-> **NOTE:** `approval_required` can only be set when `subscription_required` is set to `true`.

* `display_name` - (Required) The Display Name for this API Management Product.

* `product_id` - (Required) The Identifier for this Product, which must be unique within the API Management Service. Changing this forces a new resource to be created.

* `published` - (Required) Is this Product Published?

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service should be exist. Changing this forces a new resource to be created.

* `subscription_required` - (Required) Is a Subscription required to access API's included in this Product?

---

* `description` - (Optional) A description of this Product, which may include HTML formatting tags.

* `subscriptions_limit` - (Optional) The number of subscriptions a user can have to this Product at the same time.

-> **NOTE:** `subscriptions_limit` can only be set when `subscription_required` is set to `true`. 

* `terms` - (Optional) The Terms and Conditions for this Product, which must be accepted by Developers before they can begin the Subscription process.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Product.

## Import

API Management Products can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_product.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/products/myproduct
```
