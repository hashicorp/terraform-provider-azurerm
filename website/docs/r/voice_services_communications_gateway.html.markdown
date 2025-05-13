---
subcategory: "Voice Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_voice_services_communications_gateway"
description: |-
  Manages a Voice Services Communications Gateways.
---

# azurerm_voice_services_communications_gateway

Manages a Voice Services Communications Gateways.

!> **Note:** You must have signed an Operator Connect agreement with Microsoft to use this resource. For more information, see [`Prerequisites`](https://learn.microsoft.com/en-us/azure/communications-gateway/prepare-to-deploy#prerequisites).

!> **Note:** Access to Azure Communications Gateway is restricted, see [`Get access to Azure Communications Gateway for your Azure subscription`](https://learn.microsoft.com/en-us/azure/communications-gateway/prepare-to-deploy#9-get-access-to-azure-communications-gateway-for-your-azure-subscription) for details.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_voice_services_communications_gateway" "example" {
  name                = "example-vcg"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  connectivity        = "PublicAddress"
  codecs              = "PCMA"
  e911_type           = "DirectToEsrp"
  platforms           = ["OperatorConnect", "TeamsPhoneMobile"]

  service_location {
    location                                  = "eastus"
    allowed_media_source_address_prefixes     = ["10.1.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.1.1.0/24"]
    esrp_addresses                            = ["198.51.100.3"]
    operator_addresses                        = ["198.51.100.1"]
  }

  service_location {
    location                                  = "eastus2"
    allowed_media_source_address_prefixes     = ["10.2.2.0/24"]
    allowed_signaling_source_address_prefixes = ["10.2.1.0/24"]
    esrp_addresses                            = ["198.51.100.4"]
    operator_addresses                        = ["198.51.100.2"]
  }
  auto_generated_domain_name_label_scope = "SubscriptionReuse"
  api_bridge                             = jsonencode({})
  emergency_dial_strings                 = ["911", "933"]
  on_prem_mcp_enabled                    = false

  tags = {
    key = "value"
  }

  microsoft_teams_voicemail_pilot_number = "1"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Voice Services Communications Gateways. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Voice Services Communications Gateways should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Voice Services Communications Gateways should exist. Changing this forces a new resource to be created.

* `connectivity` - (Required) How to connect back to the operator network, e.g. MAPS. Possible values is `PublicAddress`. Changing this forces a new Voice Services Communications Gateways to be created.

* `codecs` - (Required) The voice codecs expected for communication with Teams. Possible values are `PCMA`, `PCMU`,`G722`,`G722_2`,`SILK_8` and `SILK_16`.

* `e911_type` - (Required) How to handle 911 calls. Possible values are `Standard` and `DirectToEsrp`.

* `platforms` - (Required) The Voice Services Communications GatewaysAvailable supports platform types. Possible values are `OperatorConnect`, `TeamsPhoneMobile`.

* `service_location` - (Required) A `service_location` block as defined below.

* `auto_generated_domain_name_label_scope` - (Optional) Specifies the scope at which the auto-generated domain name can be re-used. Possible values are `TenantReuse`, `SubscriptionReuse`, `ResourceGroupReuse` and `NoReuse` . Changing this forces a new resource to be created. Defaults to `TenantReuse`.

* `api_bridge` - (Optional) Details of API bridge functionality, if required.

* `emergency_dial_strings` - (Optional) A list of dial strings used for emergency calling.

* `on_prem_mcp_enabled` - (Optional) Whether an on-premises Mobile Control Point is in use.

* `tags` - (Optional) A mapping of tags which should be assigned to the Voice Services Communications Gateways.

* `microsoft_teams_voicemail_pilot_number` - (Optional) This number is used in Teams Phone Mobile scenarios for access to the voicemail IVR from the native dialer.

---

A `service_location` block supports the following:

* `location` - (Required) Specifies the region in which the resources needed for Teams Calling will be deployed.

* `operator_addresses` - (Required) IP address to use to contact the operator network from this region.

* `allowed_media_source_address_prefixes` - (Optional) Specifies the allowed source IP address or CIDR ranges for media.

* `allowed_signaling_source_address_prefixes` - (Optional) Specifies the allowed source IP address or CIDR ranges for signaling.

* `esrp_addresses` - (Optional) IP address to use to contact the ESRP from this region.

!> **Note:** The `esrp_addresses` must be specified for each `service_location` when the`e911_type` is set to `DirectToEsrp`.  The `esrp_addresses` must not be specified for each `service_location` when the`e911_type` is set to `Standard`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Voice Services Communications Gateways.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Voice Services Communications Gateways.
* `read` - (Defaults to 5 minutes) Used when retrieving the Voice Services Communications Gateways.
* `update` - (Defaults to 30 minutes) Used when updating the Voice Services Communications Gateways.
* `delete` - (Defaults to 30 minutes) Used when deleting the Voice Services Communications Gateways.

## Import

Voice Services Communications Gateways can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_voice_services_communications_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.VoiceServices/communicationsGateways/communicationsGateway1
```
