---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_api_policy"
description: |-
  Manages an API Management API Policy
---

# azurerm_api_management_api_policy

Manages an API Management API Policy

## Example Usage

```hcl
data "azurerm_api_management_api" "example" {
  name                = "my-api"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
  revision            = "2"
}

resource "azurerm_api_management_api_policy" "example" {
  api_name            = data.azurerm_api_management_api.example.name
  api_management_name = data.azurerm_api_management_api.example.api_management_name
  resource_group_name = data.azurerm_api_management_api.example.resource_group_name

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

* `api_name` - (Required) The ID of the API Management API within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `xml_content` - (Optional) The XML Content for this Policy as a string. An XML file can be used here with Terraform's [file function](https://www.terraform.io/docs/configuration/functions/file.html) that is similar to Microsoft's `PolicyFilePath` option. If you need to pass variables into your XML file, use Terraform's [templatefile function](https://developer.hashicorp.com/terraform/language/functions/templatefile).


* `xml_link` - (Optional) A link to a Policy XML Document, which must be publicly available.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management API Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management API Policy.
* `update` - (Defaults to 30 minutes) Used when updating the API Management API Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management API Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management API Policy.

## Import

API Management API Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_api_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/apis/exampleId
```
