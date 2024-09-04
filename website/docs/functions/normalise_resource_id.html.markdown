---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: normalise_resource_id"
description: |-
  Normalises a supported Azure Resource Manager ID to the correct casing for Terraform.
---

# Function: normalise_resource_id

~> Provider-defined functions are supported in Terraform 1.8 and later, and are available from version 4.0 of the provider.

~> **NOTE:** This function is also available during the opt-in beta for 4.0, available from v3.114.0. See the [beta opt-in guide](website/docs/guides/4.0-beta.html.markdown) for more information.

Takes an Azure Resource ID and attempts to normalise the case-sensitive system segments as required by the AzureRM provider. 

~> **NOTE:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result. 

## Example Usage

```hcl
# result: /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1

output "test" {
  value = provider::azurerm::normalise_resource_id("/Subscriptions/12345678-1234-9876-4563-123456789012/ResourceGroups/resGroup1/PROVIDERS/microsoft.apimanagement/service/service1/gateWays/gateway1/hostnameconfigurations/config1")
}

```

## Example - Import
```hcl
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
