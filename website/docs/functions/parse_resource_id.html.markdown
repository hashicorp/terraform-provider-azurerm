---
layout: "azurerm"
page_title: "Function: parse_resource_id"
subcategory: ""
description: |-
  Parses an Azure Resource Manager ID and exposes the contained information
  ~> Note: User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.
---

# Function: parse_resource_id

Parses an Azure Resource Manager ID and exposes the contained information

~> **Note:** User specified segments are not affected or corrected. (e.g. resource names). Please ensure that these match your configuration correctly to avoid errors. If a resource is not supported by the provider, this function may not provide a correct result.

## Example Usage

```terraform
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
parse_resource_id(id string) object
```

## Arguments

1. `id` (String) Azure Resource Manager ID.

## Return Type

The function returns an object with the following attributes:

- `full_resource_type` (String)
- `parent_resources` (Map of Strings)
- `resource_group_name` (String)
- `resource_name` (String)
- `resource_provider` (String)
- `resource_scope` (String)
- `resource_type` (String)
- `subscription_id` (String)
