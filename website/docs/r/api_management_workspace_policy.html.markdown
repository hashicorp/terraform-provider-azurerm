---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_workspace_policy"
description: |-
  Manages an API Management Workspace Policy.
---

# azurerm_api_management_workspace_policy

Manages an API Management Workspace Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apimanagement"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Premium_1"
}

resource "azurerm_api_management_workspace" "example" {
  name              = "example-workspace"
  api_management_id = azurerm_api_management.example.id
  display_name      = "my workspace"
}

resource "azurerm_api_management_workspace_policy" "example" {
  api_management_workspace_id = azurerm_api_management_workspace.example.id
  xml_content                 = <<XML
<policies>
  <inbound>
    <find-and-replace from="abc" to="xyz" />
  </inbound>
</policies>
XML
}
```

## Arguments Reference

The following arguments are supported:

* `api_management_workspace_id` - (Required) Specifies the ID of the API Management Workspace. Changing this forces a new resource to be created.

---

* `xml_content` - (Optional) Specifies the API Management Workspace Policy as an XML string.

* `xml_link` - (Optional) Specifies a publicly accessible URL to a policy XML document.

~> **Note:** Exactly one of `xml_content` or `xml_link` must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Workspace Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Workspace Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Workspace Policy.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Workspace Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Workspace Policy.

## Import

API Management Workspace Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_workspace_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/workspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.ApiManagement` - 2024-05-01
