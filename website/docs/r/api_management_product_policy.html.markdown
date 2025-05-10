---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_product_policy"
description: |-
  Manages an API Management Product Policy
---

# azurerm_api_management_product_policy

Manages an API Management Product Policy

## Example Usage

```hcl
data "azurerm_api_management_product" "example" {
  product_id          = "my-product"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
}

resource "azurerm_api_management_product_policy" "example" {
  product_id          = data.azurerm_api_management_product.example.product_id
  api_management_name = data.azurerm_api_management_product.example.api_management_name
  resource_group_name = data.azurerm_api_management_product.example.resource_group_name

  xml_content = <<XML
<policies>
  <inbound>
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML

}
```

## Argument Reference

The following arguments are supported:

* `product_id` - (Required) The ID of the API Management Product within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `xml_content` - (Optional) The XML Content for this Policy.

* `xml_link` - (Optional) A link to a Policy XML Document, which must be publicly available.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Product Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Product Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Product Policy.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Product Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Product Policy.

## Import

API Management Product Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_product_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/products/product1
```
