---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_devspace_controller"
sidebar_current: "docs-azurerm-resource-devspace-controller"
description: |-
  Manages a DevSpace Controller.
---

# azurerm_devspace_controller

Manages a DevSpace Controller.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG1"
  location = "westeurope"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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

resource "azurerm_devspace_controller" "test" {
  name                = "acctestdsc1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "S1"
    tier = "Standard"
  }

  host_suffix                              = "suffix"
  target_container_host_resource_id        = "${azurerm_kubernetes_cluster.test.id}"
  target_container_host_credentials_base64 = "${base64encode(azurerm_kubernetes_cluster.test.kube_config_raw)}"

  tags = {
    Environment = "Testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the DevSpace Controller. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the DevSpace Controller resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported location where the DevSpace Controller should exist. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as documented below. Changing this forces a new resource to be created.

* `host_suffix` - (Required) The host suffix for the DevSpace Controller. Changing this forces a new resource to be created.

* `target_container_host_resource_id` - (Required) The resource id of Azure Kubernetes Service cluster. Changing this forces a new resource to be created.

* `target_container_host_credentials_base64` - (Required) Base64 encoding of `kube_config_raw` of Azure Kubernetes Service cluster. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU for DevSpace Controller. Currently the only supported value is `S1`. Changing this forces a new resource to be created.
* `tier` - (Required) The tier of the SKU for DevSpace Controller. Currently the only supported value is `Standard`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DevSpace Controller.

* `data_plane_fqdn` - DNS name for accessing DataPlane services.

## Import

DevSpace Controller's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_devspace_controller.controller1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevSpaces/controllers/controller1Name
```
