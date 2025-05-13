---
subcategory: "Security Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iot_security_solution"
description: |-
  Manages an iot security solution.
---

# azurerm_iot_security_solution

Manages an iot security solution.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_iothub" "example" {
  name                = "example-IoTHub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iot_security_solution" "example" {
  name                = "example-Iot-Security-Solution"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  display_name        = "Iot Security Solution"
  iothub_ids          = [azurerm_iothub.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Iot Security Solution. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which to create the Iot Security Solution. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `display_name` - (Required) Specifies the Display Name for this Iot Security Solution.

* `iothub_ids` - (Required) Specifies the IoT Hub resource IDs to which this Iot Security Solution is applied.

* `additional_workspace` - (Optional) A `additional_workspace` block as defined below.

* `disabled_data_sources` - (Optional) A list of disabled data sources for the Iot Security Solution. Possible value is `TwinData`.

* `enabled` - (Optional) Is the Iot Security Solution enabled? Defaults to `true`.

* `events_to_export` - (Optional) A list of data which is to exported to analytic workspace. Valid values include `RawEvents`.

* `log_analytics_workspace_id` - (Optional) Specifies the Log Analytics Workspace ID to which the security data will be sent.

* `log_unmasked_ips_enabled` - (Optional) Should IP addressed be unmasked in the log? Defaults to `false`.

* `recommendations_enabled` - (Optional) A `recommendations_enabled` block of options to enable or disable as defined below.

* `query_for_resources` - (Optional) An Azure Resource Graph query used to set the resources monitored.

* `query_subscription_ids` - (Optional) A list of subscription Ids on which the user defined resources query should be executed.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `additional_workspace` block supports the following:

* `data_types` - (Required) A list of data types which sent to workspace. Possible values are `Alerts` and `RawEvents`.

* `workspace_id` - (Required) The resource ID of the Log Analytics Workspace.

---

A `recommendations_enabled` block supports the following:

* `acr_authentication` - (Optional) Is Principal Authentication enabled for the ACR repository? Defaults to `true`.

* `agent_send_unutilized_msg` - (Optional) Is Agent send underutilized messages enabled? Defaults to `true`.

* `baseline` - (Optional) Is Security related system configuration issues identified? Defaults to `true`.

* `edge_hub_mem_optimize` - (Optional) Is IoT Edge Hub memory optimized? Defaults to `true`.

* `edge_logging_option` - (Optional) Is logging configured for IoT Edge module? Defaults to `true`.

* `inconsistent_module_settings` - (Optional) Is inconsistent module settings enabled for SecurityGroup? Defaults to `true`.

* `install_agent` - (Optional) is Azure IoT Security agent installed? Defaults to `true`.

* `ip_filter_deny_all` - (Optional) Is Default IP filter policy denied? Defaults to `true`.

* `ip_filter_permissive_rule` - (Optional) Is IP filter rule source allowable IP range too large? Defaults to `true`.

* `open_ports` - (Optional) Is any ports open on the device? Defaults to `true`.

* `permissive_firewall_policy` - (Optional) Does firewall policy exist which allow necessary communication to/from the device? Defaults to `true`.

* `permissive_input_firewall_rules` - (Optional) Is only necessary addresses or ports are permitted in? Defaults to `true`.

* `permissive_output_firewall_rules` - (Optional) Is only necessary addresses or ports are permitted out? Defaults to `true`.

* `privileged_docker_options` - (Optional) Is high level permissions are needed for the module? Defaults to `true`.

* `shared_credentials` - (Optional) Is any credentials shared among devices? Defaults to `true`.

* `vulnerable_tls_cipher_suite` - (Optional) Does TLS cipher suite need to be updated? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Iot Security Solution resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Iot Security Solution.
* `read` - (Defaults to 5 minutes) Used when retrieving the Iot Security Solution.
* `update` - (Defaults to 30 minutes) Used when updating the Iot Security Solution.
* `delete` - (Defaults to 30 minutes) Used when deleting the Iot Security Solution.

## Import

Iot Security Solution can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_security_solution.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Security/iotSecuritySolutions/solution1
```
