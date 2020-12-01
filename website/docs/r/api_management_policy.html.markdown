---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_policy"
description: |-
  Manages a API Management service Policy.
---

# azurerm_api_management_policy

Manages a API Management service Policy.

~> **NOTE:** This resource will, upon creation, **overwrite any existing policy in the API Management service**, as there is no feasible way to test whether the policy has been modified from the default. Similarly, when this resource is destroyed, the API Management service will revert to its default policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_named_value" "example" {
  name                = "example-apimg"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  display_name        = "ExampleProperty"
  value               = "Example Value"
}

resource "azurerm_api_management_policy" "example" {
  apim_management_id = azurerm_api_management.example.id
  xml_content        = file("example.xml")
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_id` - (Required) The ID of the API Management service. Changing this forces a new API Management service Policy to be created.

---

* `xml_content` - (Optional) The XML Content for this Policy as a string. An XML file can be used here with Terraform's [file function](https://www.terraform.io/docs/configuration/functions/file.html) that is similar to Microsoft's `PolicyFilePath` option.

* `xml_link` - (Optional) A link to a Policy XML Document, which must be publicly available.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management service Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management service Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management service Policy.
* `update` - (Defaults to 30 minutes) Used when updating the API Management service Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management service Policy.

## Import

API Management service Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/policies/policy
```
