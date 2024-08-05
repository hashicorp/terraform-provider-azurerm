---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_cluster"
description: |-
  Gets information about an existing Azure Stack HCI Cluster.
---

# Data Source: azurerm_stack_hci_cluster

Use this data source to access information about an existing Azure Stack HCI Cluster instance.

## Example Usage

```hcl
data "azurerm_stack_hci_cluster" "example" {
  name                = "existing"
  resource_group_name = "existing"
}


output "id" {
  value = data.azurerm_stack_hci_cluster.example.id
}

output "location" {
  value = data.azurerm_stack_hci_cluster.example.location
}

output "client_id" {
  value = data.azurerm_stack_hci_cluster.example.client_id
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Stack HCI Cluster.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Cluster exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Stack HCI Cluster.

* `automanage_configuration_id` - The ID of the Automanage Configuration assigned to the Azure Stack HCI Cluster.

* `client_id` - The Client ID of the Azure Active Directory used by the Azure Stack HCI Cluster.

* `cloud_id` - An immutable UUID for the Azure Stack HCI Cluster.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Azure Stack HCI Cluster exists.

* `resource_provider_object_id` - The object ID of the Resource Provider Service Principal.

* `service_endpoint` - The region specific Data Path Endpoint of the Azure Stack HCI Cluster.

* `tenant_id` - The Tenant ID of the Azure Active Directory used by the Azure Stack HCI Cluster.

* `tags` - A mapping of tags assigned to the Azure Stack HCI Cluster.

---

An `identity` block exports the following:

* `type` - (Required) The type of Managed Service Identity configured on the Azure Stack HCI Cluster.

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Cluster.
