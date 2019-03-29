layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_policy"
sidebar_current: "docs-azurerm-resource-api-management-policy"
description: |-
  Manages a global Policy within an API Management Service.
---

# azurerm_api_management_policy

Manages a global Policy within an API Management Service.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West US"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku = {
    name     = "Developer"
    capacity = 1
  }
}

resource "azurerm_api_management_policy" "example" {
  resource_group_name = "${azurerm_resource_group.example.name}"
  api_management_name = "${azurerm_api_management.example.name}"
  xml_content         = "<policies><inbound></inbound></policies>"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service where this Policy should be created. Changing this forces a new resource to be created.

* `xml_content` - (Optional) The XML configuration of this Policy.

* `xml_link` - (Optional) The HTTP endpoint of the XML configuration accessible from the API Management Service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the API Management Policy.


## Import

API Management Policy can be imported using the `resource id`, e.g.
```shell
$ terraform import azurerm_api_management_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/Microsoft.ApiManagement/service/example-apim/policies/policy
```
