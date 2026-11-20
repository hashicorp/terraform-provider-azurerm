---
layout: "azurerm"
page_title: "Function: normalise_resource_id"
subcategory: ""
description: |-
  Takes an Azure Resource ID and attempts to normalise the case-sensitive system segments as required by the AzureRM provider.
  ~> Note: User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.
---

# Function: normalise_resource_id

Takes an Azure Resource ID and attempts to normalise the case-sensitive system segments as required by the AzureRM provider.

~> **Note:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.

## Example Usage

```terraform
# result: /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1

output "test" {
  value = provider::azurerm::normalise_resource_id("/Subscriptions/12345678-1234-9876-4563-123456789012/ResourceGroups/resGroup1/PROVIDERS/microsoft.apimanagement/service/service1/gateWays/gateway1/hostnameconfigurations/config1")
}
```

## Example - Import

```terraform
import {
  id = provider::azurerm::normalise_resource_id("/Subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/import-example")
  to = azurerm_resource_group.test
}

resource "azurerm_resource_group" "test" {
  name     = "import-example"
  location = "westeurope"
}
```

## Signature

```text
normalise_resource_id(id string) string
```

## Arguments

1. `id` (String) Azure Resource Manager ID.
