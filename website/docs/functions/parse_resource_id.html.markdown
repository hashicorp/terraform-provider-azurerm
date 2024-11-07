---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: normalise_resource_id"
description: |-
  Normalises a supported Azure Resource Manager ID to the correct casing for Terraform.
---

# Function: parse_resource_id

~> Provider-defined functions are supported in Terraform 1.8 and later, and are available from version 4.0 of the provider.

~> **NOTE:** This function is also available during the opt-in beta for 4.0, available from v3.114.0. See the [beta opt-in guide](website/docs/guides/4.0-beta.html.markdown) for more information.

Takes an Azure Resource ID and splits it into its component parts. 

~> **NOTE:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result. 

## Example Usage

```hcl
# result:
# Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
# 
# Outputs:
# 
# parsed_id = {
# "full_resource_type" = "Microsoft.ApiManagement/service/gateways/hostnameConfigurations"
# "parent_resources" = tomap({
# "gateways" = "gateway1"
# "service" = "service1"
# })
# "resource_group_name" = "resGroup1"
# "resource_name" = "config1"
# "resource_provider" = "Microsoft.ApiManagement"
# "resource_scope" = tostring(null)
# "resource_type" = "hostnameConfigurations"
# "subscription_id" = "12345678-1234-9876-4563-123456789012"
# }
# resource_name = "config1"

provider "azurerm" {
  features {}
}

locals {
  parsed_id = provider::azurerm::parse_resource_id("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1")
}

output "parsed" {
  value = local.parsed_id
}

output "resource_name" {
  value = local.parsed_id["resource_name"]
}


```

## Signature

```text
parse_resource_id(id string) string
```

## Arguments

1. `id` (String) Azure Resource Manager ID.
