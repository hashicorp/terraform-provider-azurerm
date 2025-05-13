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
  display_name = "Allowed resource types"
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
  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Cluster should exist. Changing this forces a new resource to be created.

* `client_id` - (Optional) The Client ID of the Azure Active Directory Application which is used by the Azure Stack HCI Cluster. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tenant_id` - (Optional) The Tenant ID of the Azure Active Directory which is used by the Azure Stack HCI Cluster. Changing this forces a new resource to be created.

~> **Note:** If unspecified the Tenant ID of the Provider will be used.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Cluster.

* `automanage_configuration_id` - (Optional) The ID of the Automanage Configuration assigned to the Azure Stack HCI Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on the Azure Stack HCI Cluster. Possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Azure Stack HCI Cluster.

* `cloud_id` - An immutable UUID for the Azure Stack HCI Cluster.

* `resource_provider_object_id` - The object ID of the Resource Provider Service Principal.

* `identity` - An `identity` block as defined below.

* `service_endpoint` - The region specific Data Path Endpoint of the Azure Stack HCI Cluster.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

-> **Note:** You can access the Principal ID via `azurerm_stack_hci_cluster.example.identity.0.principal_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Cluster.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Cluster.

## Import

Azure Stack HCI Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1
```
