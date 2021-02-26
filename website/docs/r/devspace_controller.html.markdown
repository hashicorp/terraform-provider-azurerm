---
subcategory: "DevSpace"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_devspace_controller"
description: |-
  Manages a DevSpace Controller.
---

# azurerm_devspace_controller

Manages a DevSpace Controller.

~> **NOTE:** Microsoft will be retiring Azure Dev Spaces on 31 October 2023, please see the product [documentation](https://azure.microsoft.com/en-us/updates/azure-dev-spaces-is-retiring-on-31-october-2023/) for more information.

!> **NOTE:** The Azure API no longer allows provisioning new DevSpace Controllers - as such this resource exists only to allow existing users to continue managing these in Terraform at this time. Support for the `azurerm_devspace_controller` resource will be removed in version 3.0 of the Azure Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example_resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "acctestaks1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "acctestaks1"

  agent_pool_profile {
    name    = "default"
    count   = "1"
    vm_size = "Standard_DS2_v2"
  }

  service_principal {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "00000000000000000000000000000000"
  }
}

resource "azurerm_devspace_controller" "example" {
  name                = "acctestdsc1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "S1"

  host_suffix                              = "suffix"
  target_container_host_resource_id        = azurerm_kubernetes_cluster.example.id
  target_container_host_credentials_base64 = "${base64encode(azurerm_kubernetes_cluster.example.kube_config_raw)}"

  tags = {
    Environment = "Testing"
  }
}
```

## Argument Reference

!> **NOTE:** The Azure API no longer allows provisioning new DevSpace Controllers - as such this resource exists only to allow existing users to continue managing these in Terraform at this time. Support for the `azurerm_devspace_controller` resource will be removed in version 3.0 of the Azure Provider.

The following arguments are supported:  

* `name` - (Required) Specifies the name of the DevSpace Controller. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the DevSpace Controller resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported location where the DevSpace Controller should exist. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this DevSpace Controller. Possible values are `S1`.

* `target_container_host_resource_id` - (Required) The resource id of Azure Kubernetes Service cluster. Changing this forces a new resource to be created.

* `target_container_host_credentials_base64` - (Required) Base64 encoding of `kube_config_raw` of Azure Kubernetes Service cluster. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DevSpace Controller.

* `data_plane_fqdn` - DNS name for accessing DataPlane services.

* `host_suffix` - The host suffix for the DevSpace Controller.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevSpace Controller.
* `update` - (Defaults to 30 minutes) Used when updating the DevSpace Controller.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevSpace Controller.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevSpace Controller.

## Import

DevSpace Controller's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_devspace_controller.controller1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevSpaces/controllers/controller1Name
```
