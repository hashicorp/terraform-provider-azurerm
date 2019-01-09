---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
sidebar_current: "docs-azurerm-resource-container-group"
description: |-
  Create as an Azure Container Group instance.
---

# azurerm_container_group

Manage as an Azure Container Group instance.

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
  os_type             = "Linux"

  container {
    name   = "hw"
    image  = "seanmckenna/aci-hellofiles"
    cpu    = "0.5"
    memory = "1.5"
    ports  = {
      port = 80
    }
    ports = {
      port = 443 
    }

    environment_variables {
      "NODE_ENV" = "testing"
    }

    secure_environment_variables {
      "ACCESS_KEY" = "secure_testing"
    }

    commands = ["/bin/bash", "-c", "'/path to/myscript.sh'"]

    volume {
      name       = "logs"
      mount_path = "/aci/logs"
      read_only  = false
      share_name = "${azurerm_storage_share.aci-share.name}"

      storage_account_name = "${azurerm_storage_account.aci-sa.name}"
      storage_account_key  = "${azurerm_storage_account.aci-sa.primary_access_key}"
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

* `image_registry_credential` - (Optional) Set image registry credentials for the group as documented in the `image_registry_credential` block below

* `container` - (Required) The definition of a container that is part of the group as documented in the `container` block below. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `ports` - (Optional) A set of public ports for the container. Changing this forces a new resource to be created. Set as documented in the `ports` block below.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `secure_environment_variables` - (Optional) A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container.

~> **NOTE:** The field `command` has been deprecated in favor of `commands` to better match the API.

* `commands` - (Optional) A list of commands which should be run on the container.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

The `volume` block supports:

* `name` - (Required) The name of the volume mount. Changing this forces a new resource to be created.

* `mount_path` - (Required) The path on which this volume is to be mounted. Changing this forces a new resource to be created.

* `read_only` - (Optional) Specify if the volume is to be mounted as read only or not. The default value is `false`. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The Azure storage account from which the volume is to be mounted. Changing this forces a new resource to be created.

* `storage_account_key` - (Required) The access key for the Azure Storage account specified as above. Changing this forces a new resource to be created.

* `share_name` - (Required) The Azure storage share that is to be mounted as a volume. This must be created on the storage account specified as above. Changing this forces a new resource to be created.

The `image_registry_credential` block supports:

* `username` - (Required) The username with which to connect to the registry.

* `password` - (Required) The password with which to connect to the registry.

* `server` - (Required) The address to use to connect to the registry without protocol ("https"/"http"). For example: "myacr.acr.io"

The `ports` block supports:

* `port` - (Required) The port number the container will expose.

* `protocol` - (Optional) The network protocol ("tcp"/"udp") associated with port. The default is `TCP`.

## Attributes Reference

The following attributes are exported:

* `id` - The container group ID.

* `ip_address` - The IP address allocated to the container group.

* `fqdn` - The FQDN of the container group derived from `dns_name_label`.

## Import

Container Group's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_group.containerGroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerInstance/containerGroups/myContainerGroup1
```
