---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_diagnostic"
sidebar_current: "docs-azurerm-resource-api-management-diagnostic"
description: |-
  Manages an API Management Service Diagnostic.
---

# azurerm_api_management_diagnostic

Manages an API Management Service Diagnostic.

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
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_diagnostic" "test" {
  identifier          = "applicationinsights"
  resource_group_name = "${azurerm_resource_group.test.name}"
  api_management_name = "${azurerm_api_management.test.name}"
  enabled             = true
}
```

## Argument Reference

The following arguments are supported:

* `identifier` - (Required) The diagnostic identifier for the API Management Service. At this time the only supported value is `applicationinsights`. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The Name of the API Management Service where this Diagnostic should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management Service exists. Changing this forces a new resource to be created.

* `enabled` - (Required) Indicates whether a Diagnostic should receive data or not.

---

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Diagnostic.

## Import

API Management Diagnostics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_diagnostic.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/diagnostics/applicationinsights
```
