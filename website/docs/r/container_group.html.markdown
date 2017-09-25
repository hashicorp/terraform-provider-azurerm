---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
sidebar_current: "docs-azurerm-resource-container-group"
description: |-
  Create as an Azure Container Group instance.
---

# azurerm\_container\_group

Create as an Azure Container Group instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "aci-rg" {
  name     = "aci-test"
  location = "west us"
}

resource "azurerm_container_group" "aci-helloworld" {
  name                = "aci-hw"
  location            = "${azurerm_resource_group.aci-rg.location}"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type     = "public"
  os_type             = "linux"

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"
    port   = "80"

    environment_variables {
      "NODE_ENV" = "Staging"
    }

    command = "/bin/bash -c '/path to/myscript.sh'"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
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

* `ip_address_type` - (Optional) Specifies the ip address type of the container. `Public` is the only acceptable value at this time. Changing this forces a new resource to be created.

* `os_type` - (Required) The OS for the container group. Allowed values are `Linux` and `Windows` Changing this forces a new resource to be created.

* `container` - (Required) The definition of a container that is part of the group. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported.

The `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `port` - (Optional) A public port for the container. Changing this forces a new resource to be created.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The container group ID.

* `ip_address` - The IP address allocated to the container group.
