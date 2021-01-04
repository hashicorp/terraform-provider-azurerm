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

* `iothub_ids` - (Required) Specifies the IoT Hub resource IDs to which this Iot Security Solution applies.

* `log_analytics_workspace_id` - (Optional) Specifies the Log Analytics Workspace ID to which the security data will be sent.

* `enabled` - (Optional) Is the Iot Security Solution enabled? Defaults to `true`.

* `unmasked_ip_logging_enabled` - (Optional) Is unmasked IP logging enabled? Defaults to `false`.

* `export` - (Optional) A list of data which is to exported to analytic workspace. Valid values include `RawEvents`.

* `recommendation` - (Optional) A `recommendation` block as defined below.

* `user_defined_resource` - (Optional) A `user_defined_resource` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `recommendation` block supports the following:

* `iot_acr_authentication_enabled` - (Optional) Could Service Principal Authentication be used with ACR repository? Defaults to `true`.

* `iot_agent_send_unutilized_msg_enabled` - (Optional) Could Agent send underutilized messages? Defaults to `true`.

* `iot_baseline_enabled` - (Optional) Is Security related system configuration issues identified? Defaults to `true`.

* `iot_edge_hub_mem_optimize_enabled` - (Optional) Is IoT Edge Hub memory optimized? Defaults to `true`.

* `iot_edge_logging_option_enabled` - (Optional) Is logging configured for IoT Edge module? Defaults to `true`.

* `iot_inconsistent_module_settings_enabled` - (Optional) Does SecurityGroup has inconsistent module settings? Defaults to `true`.

* `iot_install_agent_enabled` - (Optional) is Azure IoT Security agent installed? Defaults to `true`.

* `iot_ip_filter_deny_all_enabled` - (Optional) Is Default IP filter policy denied? Defaults to `true`.

* `iot_ip_filter_permissive_rule_enabled` - (Optional) Is IP filter rule source allowable IP range too large? Defaults to `true`.

* `iot_open_ports_enabled` - (Optional) Is any ports open on the device? Defaults to `true`.

* `iot_permissive_firewall_policy_enabled` - (Optional) Does firewall policy exist which allow necessary communication to/from the device? Defaults to `true`.

* `iot_permissive_input_firewall_rules_enabled` - (Optional) Is only necessary addresses or ports are permitted in? Defaults to `true`.

* `iot_permissive_output_firewall_rules_enabled` - (Optional) Is only necessary addresses or ports are permitted out? Defaults to `true`.

* `iot_privileged_docker_options_enabled` - (Optional) Is high level permissions are needed for the module? Defaults to `true`.

* `iot_shared_credentials_enabled` - (Optional) Is any credentials shared among devices? Defaults to `true`.

* `iot_vulnerable_tls_cipher_suite_enabled` - (Optional) Does TLS cipher suite need to be updated? Defaults to `true`.

---

A `user_defined_resource` block supports the following:

* `query` - (Required) Azure Resource Graph query which represents the security solution's user defined resources.

* `subscription_ids` - (Required) A list of subscription Ids on which the user defined resources query should be executed..

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Advanced Threat Protection resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Iot Security Solution.
* `update` - (Defaults to 30 minutes) Used when updating the Iot Security Solution.
* `read` - (Defaults to 5 minutes) Used when retrieving the Iot Security Solution.
* `delete` - (Defaults to 30 minutes) Used when deleting the Iot Security Solution.

## Import

Iot Security Solution can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iot_security_solution.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Security/IoTSecuritySolutions/solution1
```
