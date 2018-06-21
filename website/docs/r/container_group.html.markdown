---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
sidebar_current: "docs-azurerm-resource-container-group"
description: |-
  Create as an Azure Container Group instance.
---

# azurerm_container_group

Create as an Azure Container Group instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "aci-rg" {
  name     = "aci-test"
  location = "west us"
}

resource "azurerm_storage_account" "aci-sa" {
  name                = "acistorageacct"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  location            = "${azurerm_resource_group.aci-rg.location}"
  account_tier        = "Standard"
  
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "aci-share" {
  name = "aci-test-share"

  resource_group_name  = "${azurerm_resource_group.aci-rg.name}"
  storage_account_name = "${azurerm_storage_account.aci-sa.name}"

  quota = 50
}

resource "azurerm_container_group" "aci-helloworld" {
  name                = "aci-hw"
  location            = "${azurerm_resource_group.aci-rg.location}"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type     = "public"
  dns_name_label      = "aci-label"
  os_type             = "linux"

  container {
    name   = "hw"
    image  = "seanmckenna/aci-hellofiles"
    cpu    ="0.5"
    memory =  "1.5"
    port   = "80"

    environment_variables {
      "NODE_ENV" = "testing"
    }

    command = "/bin/bash -c '/path to/myscript.sh'"

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false
      share_name = "${azurerm_storage_share.aci-share.name}"
      
      storage_account_name  = "${azurerm_storage_account.aci-sa.name}"
      storage_account_key   = "${azurerm_storage_account.aci-sa.primary_access_key}"
    }
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

* `dns_name_label` - (Optional) The DNS label/name for the container groups IP.

* `os_type` - (Required) The OS for the container group. Allowed values are `Linux` and `Windows`. Changing this forces a new resource to be created.

* `restart_policy` - (Optional) Restart policy for the container group. Allowed values are `Always`, `Never`, `OnFailure`. Defaults to `Always`.

* `container` - (Required) The definition of a container that is part of the group as documented in the `container` block below. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported.

The `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `port` - (Optional) A public port for the container. Changing this forces a new resource to be created.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container. Changing this forces a new resource to be created.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

The `volume` block supports:

* `name` - (Required) The name of the volume mount. Changing this forces a new resource to be created.

* `mount_path` - (Required) The path on which this volume is to be mounted. Changing this forces a new resource to be created.

* `read_only` - (Optional) Specify if the volume is to be mounted as read only or not. The default value is `false`. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The Azure storage account from which the volume is to be mounted. Changing this forces a new resource to be created.

* `storage_account_key` - (Required) The access key for the Azure Storage account specified as above. Changing this forces a new resource to be created.

* `share_name` - (Required) The Azure storage share that is to be mounted as a volume. This must be created on the storage account specified as above. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The container group ID.

* `ip_address` - The IP address allocated to the container group.

* `fqdn` - The FQDN of the container group derived from `dns_name_label`.
