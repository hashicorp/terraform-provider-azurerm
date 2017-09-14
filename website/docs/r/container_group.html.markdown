---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
sidebar_current: "docs-azurerm-resource-container-group"
description: |-
  Create as an Azure Container Group instance.
---

# azurerm\_container\_group

Create as an Azure Container Group instance.

## Example Usage (Linux)

```hcl
resource "azurerm_resource_group" "aci-rg" {
  name     = "aci-test"
  location = "west us"
}

resource "azurerm_container_group" "aci-helloworld" {
  
  name = "aci-hw"
  location = "west us"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type="public"
  os_type = "linux"

  container {
        name = "hw"
        image = "microsoft/aci-helloworld:latest"
        cpu ="0.5"
        memory =  "1.5"
        port = "80"
    }
    container {
        name = "sidecar"
        image = "microsoft/aci-tutorial-sidecar"
        cpu="0.5"
        memory="1.5"
    }

  tags {
    environment = "testing"
  }
}
```

## Example Usage (Windows)

```hcl
resource "azurerm_resource_group" "aci-rg" {
  name     = "aci-test"
  location = "west us"
}

resource "azurerm_container_group" "winapp" {
  
  name = "mywinapp"
  location = "west us"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type="public"
  os_type = "windows"

  container {
    name = "winapp1"
    image = "winappimage:latest"
	cpu ="2.0"
    memory = "3.5"
    port = "80"
  }

  tags {
    environment = "testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Group. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `ip_address_type` - (Optional) Specifies the ip address type of the container. `Public` is the only acceptable value at this time.

* `os_type` - (Required) The OS for the container group. Allowed values are `linux` and `windows` Changing this forces a new resource to be created.

* `container` - (Required) The definition of a container that is part of the group. Currently, only single containers are supported for Windows OS and multiple containers are supported for Linux OS.

The `container` block supports:

* `name` - (Required) Specifies the name of the Container.

* `image` - (Required) The container image name.

* `cpu` - (Required) The required number of CPU cores of the containers.

* `memory` - (Required) The required memory of the containers in GB.

* `port` - (Optional) A public port for the container.

## Attributes Reference

The following attributes are exported:

* `id` - The container group ID.

* `ip_address` - The IP address allocated to the container group.
