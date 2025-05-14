---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_dynamics_application_performance_monitoring"
description: |-
  Manages a Spring Cloud Application Performance Monitoring resource for App Dynamics.
---

# azurerm_spring_cloud_app_dynamics_application_performance_monitoring

-> **Note:** This resource is only applicable for Spring Cloud Service enterprise tier

Manages a Spring Cloud Application Performance Monitoring resource for App Dynamics.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app_dynamics_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_app_dynamics_application_performance_monitoring" "example" {
  name                     = "example"
  spring_cloud_service_id  = azurerm_spring_cloud_service.example.id
  agent_account_name       = "example-agent-account-name"
  agent_account_access_key = "example-agent-account-access-key"
  controller_host_name     = "example-controller-host-name"
  agent_application_name   = "example-agent-application-name"
  agent_tier_name          = "example-agent-tier-name"
  agent_node_name          = "example-agent-node-name"
  agent_unique_host_id     = "example-agent-unique-host-id"
  controller_ssl_enabled   = true
  controller_port          = 8080
  globally_enabled         = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Spring Cloud Application Performance Monitoring resource for App Dynamics. Changing this forces a new resource to be created.

* `spring_cloud_service_id` - (Required) The ID of the Spring Cloud Service. Changing this forces a new resource to be created.

* `agent_account_name` - (Required) Specifies the account name of the App Dynamics account.

* `agent_account_access_key` - (Required) Specifies the account access key used to authenticate with the Controller.

* `controller_host_name` - (Required) Specifies the hostname or the IP address of the AppDynamics Controller.

---

* `agent_application_name` - (Optional) Specifies the name of the logical business application that this JVM node belongs to.

* `agent_tier_name` - (Optional) Specifies the name of the tier that this JVM node belongs to.

* `agent_node_name` - (Optional) Specifies the name of the node. Where JVMs are dynamically created.

* `agent_unique_host_id` - (Optional) Specifies the unique host ID which is used to Logically partition a single physical host or virtual machine such that it appears to the Controller that the application is running on different machines.

* `controller_ssl_enabled` - (Optional) Specifies whether enable use SSL (HTTPS) to connect to the AppDynamics Controller.

* `controller_port` - (Optional) Specifies the HTTP(S) port of the AppDynamics Controller. This is the port used to access the AppDynamics browser-based user interface.

* `globally_enabled` - (Optional) Specifies whether the Spring Cloud Application Performance Monitoring resource for Application Insights is enabled globally. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application Performance Monitoring resource for App Dynamics.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application Performance Monitoring resource for App Dynamics.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application Performance Monitoring resource for App Dynamics.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application Performance Monitoring resource for App Dynamics.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application Performance Monitoring resource for App Dynamics.

## Import

Spring Cloud Application Performance Monitoring resource for App Dynamics can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app_dynamics_application_performance_monitoring.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.AppPlatform/spring/service1/apms/apm1
```
