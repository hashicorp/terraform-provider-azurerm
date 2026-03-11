---
subcategory: "Voice Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_voice_services_communications_gateway_test_line"
description: |-
  Lists Voice Services Communications Gateway Test Line resources.
---

# List resource: azurerm_voice_services_communications_gateway_test_line

Lists Voice Services Communications Gateway Test Line resources.

## Example Usage

### List all Test Lines under a Communications Gateway

```hcl
list "azurerm_voice_services_communications_gateway_test_line" "example" {
  provider = azurerm
  config {
    voice_services_communications_gateway_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.VoiceServices/communicationsGateways/example-gateway"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `voice_services_communications_gateway_id` - (Required) The ID of the Voice Services Communications Gateway whose Test Lines should be listed.
