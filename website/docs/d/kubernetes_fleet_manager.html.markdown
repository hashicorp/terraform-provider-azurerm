---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_kubernetes_fleet_manager"
description: |-
  Gets information about an existing Kubernetes Fleet Manager.
---

# Data Source: azurerm_kubernetes_fleet_manager

Use this data source to access information about an existing Kubernetes Fleet Manager.

## Example Usage

```hcl
data "azurerm_kubernetes_fleet_manager" "example" {
  name                = "example"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_kubernetes_fleet_manager.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Kubernetes Fleet Manager.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Fleet Manager exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Kubernetes Fleet Manager.

* `hub_profile` - A `hub_profile` block as defined below.

* `location` - The Azure Region where the Kubernetes Fleet Manager exists.

* `tags` - A mapping of tags assigned to the Kubernetes Fleet Manager.

---

An `agent_profile` block exports the following:

* `subnet_id` - The ID of the subnet which the Fleet hub node will join on startup. This is empty when no custom subnet is configured.

* `virtual_machine_size` - The virtual machine size of the Fleet hub.

---

An `api_server_access_profile` block exports the following:

* `enable_private_cluster` - Whether the Fleet hub is a private cluster.

---

A `hub_profile` block exports the following:

* `agent_profile` - An `agent_profile` block as defined above.

* `api_server_access_profile` - An `api_server_access_profile` block as defined above.

* `dns_prefix` - DNS prefix used to create the FQDN for the Fleet hub. This can be empty when Azure generates an FQDN without a configured DNS prefix.

* `fqdn` - The FQDN of the Fleet hub.

* `kubernetes_version` - The Kubernetes version of the Fleet hub.

* `portal_fqdn` - The Azure Portal FQDN of the Fleet hub.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Fleet Manager.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.ContainerService` - 2024-04-01
