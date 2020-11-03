---
subcategory: "AzureStackHCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_hci_cluster"
description: |-
  Manages a Azure Stack HCI Cluster.
---

# azurerm_hci_cluster

Manages a Azure Stack HCI Cluster.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_hci_cluster" "example" {
  name                = "example-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  aad_client_id       = data.azurerm_client_config.current.client_id
  aad_tenant_id       = data.azurerm_client_config.current.tenant_id

  tags = {
    ENV = "Prod"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this HCI Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the HCI Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the HCI Cluster should exist. Changing this forces a new resource to be created.

* `aad_client_id` - (Required) The ID of the AAD client. Changing this forces a new resource to be created.

* `aad_tenant_id` - (Required) The ID of the AAD tenant. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the HCI Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the HCI Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the HCI Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the HCI Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the HCI Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the HCI Cluster.

## Import

HCI Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_hci_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1
```
