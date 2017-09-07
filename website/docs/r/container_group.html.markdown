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

resource "azurerm_container_group" "nginx" {
  
  name = "mynginx"
  location = "west us"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type="public"
  os_type = "linux"

  container {
    name = "nginx1"
    image = "nginx:latest"
    cpu ="1"
    memory = "1.5"
    port = "80"
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

resource "azurerm_container_group" "nginx" {
  
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

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created. Changing this forces a new resource to be created.

* `ip_address_type` - (Required) Specifies the ip address type of the container. `public` is the only acceptable value at this time. Changing this forces a new resource to be created.

* `os_type` - (Required) The OS for the container group. Allowed values are `linux` and `windows` Changing this forces a new resource to be created.

* `container` - (Required) The definition of a container that is part of the group. Currently, only single containers are supported. Changing this forces a new resource to be created.

The `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `port` - (Optional) A public port for the container. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `ip_address` - The IP address allocated to the container group.
