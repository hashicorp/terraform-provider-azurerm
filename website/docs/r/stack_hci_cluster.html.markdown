---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_cluster"
description: |-
  Manages an Azure Stack HCI Cluster.
---

# azurerm_stack_hci_cluster

Manages an Azure Stack HCI Cluster.

## Example Usage

```hcl
data "azuread_application" "example" {
  name = "example-app"
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_stack_hci_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  client_id           = data.azuread_application.example.application_id
  tenant_id           = data.azurerm_client_config.current.tenant_id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Cluster should exist. Changing this forces a new resource to be created.

* `client_id` - (Required) The Client ID of the Azure Active Directory which is used by the Azure Stack HCI Cluster. Changing this forces a new resource to be created.

* `tenant_id` - (Optional) The Tenant ID of the Azure Active Directory which is used by the Azure Stack HCI Cluster. Changing this forces a new resource to be created.

~> **NOTE** If unspecified the Tenant ID of the Provider will be used.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Azure Stack HCI Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Cluster.

## Import

Azure Stack HCI Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1
```
