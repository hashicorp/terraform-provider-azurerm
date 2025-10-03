---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_policy_fragment"
description: |-
  Manages an API Management Workspace Policy Fragment.
---

# azurerm_api_management_workspace_policy_fragment

Manages an API Management Workspace Policy Fragment.

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

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "Example Workspace"
  description       = "Example API Management Workspace"
}

resource "azurerm_api_management_workspace_policy_fragment" "example" {
  name                        = "example-policy-fragment"
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  xml_format                  = "xml"
  xml_content                 = file("policy-fragment-1.xml")
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API Management Workspace Policy Fragment. Changing this forces a new resource to be created.

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

* `xml_content` - (Required) Specifies the XML content of the API Management Workspace Policy Fragment.

---

* `description` - (Optional) Specifies the description for the API Management Workspace Policy Fragment.

* `xml_format` - (Optional) Specifies the XML format of the API Management Workspace Policy Fragment. Possible values are `xml` or `rawxml`. Defaults to `xml`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the API Management Workspace Policy Fragment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Policy Fragment.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Policy Fragment.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Policy Fragment.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Policy Fragment.

## Import

API Management Workspace Policy Fragments can be imported using the `resource id`, e.g.

~> **Note:** Due to the behaviour of the API, API Management Workspace Policy Fragments can only be imported as `xml`, but can be updated to the desired format after importing.

```shell
terraform import azurerm_api_management_workspace_policy_fragment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/instance1/workspaces/workspace1/policyFragments/policyFragment1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
