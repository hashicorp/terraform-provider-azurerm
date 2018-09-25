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
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

#storage account name needs to be globally unique so lets generate a random id
resource "random_integer" "random_int" {
  min = 100
  max = 999
}

resource "azurerm_storage_account" "aci-sa" {
  name                     = "acistorageacct${random_integer.random_int.result}"
  resource_group_name      = "${azurerm_resource_group.aci-rg.name}"
  location                 = "${azurerm_resource_group.aci-rg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "aci-share" {
  name = "aci-test-share"

  resource_group_name  = "${azurerm_resource_group.aci-rg.name}"
  storage_account_name = "${azurerm_storage_account.aci-sa.name}"

  quota = 50
}

resource "azurerm_container_group" "aci-example" {
  name                = "mycontainergroup-${random_integer.random_int.result}"
  location            = "${azurerm_resource_group.aci-rg.location}"
  resource_group_name = "${azurerm_resource_group.aci-rg.name}"
  ip_address_type     = "public"
  dns_name_label      = "mycontainergroup-${random_integer.random_int.result}"
  os_type             = "linux"

  volume {
    name      = "emptydir"
    empty_dir = {}
  }

  volume {
    name      = "secret"
    
    secret = {
      name = "examplesecret0"
      data = "YmFzZTY0IGRhdGEK" // Base64 data saying "base64 data"
    }
    secret = {
      name = "examplesecret1"
      data = "YmFzZTY0IGRhdGEK" // Base64 data saying "base64 data"
    }
    secret = {
      name = "examplesecret2"
      data = "YmFzZTY0IGRhdGEK" // Base64 data saying "base64 data"
    }
  }


  volume {
    name = "azureshare"

    azure_share {
      share_name           = "${azurerm_storage_share.aci-share.name}"
      storage_account_name = "${azurerm_storage_account.aci-sa.name}"
      storage_account_key  = "${azurerm_storage_account.aci-sa.primary_access_key}"
    }
  }

  volume {
    name = "gitrepo"

    git_repo {
      repository = "https://github.com/Azure-Samples/aci-tutorial-sidecar"
    }
  }

  container {
    name     = "webserver"
    image    = "seanmckenna/aci-hellofiles"
    cpu      = "1"
    memory   = "1.5"
    port     = "80"
    protocol = "tcp"

    volume_mount {
      volume_name = "emptydir"
      mount_path  = "/aci/empty"
    }

    volume_mount {
      volume_name = "gitrepo"
      mount_path  = "/aci/gitrepo"
    }

    volume_mount {
      volume_name = "secret"
      mount_path  = "/aci/secret"
    }
  }

  container {
    name   = "sidecar"
    image  = "seanmckenna/aci-hellofiles"
    cpu    = "1"
    memory = "1.5"

    volume_mount {
      volume_name = "emptydir"
      mount_path  = "/empty"
      read_only   = false
    }

    volume_mount {
      volume_name = "gitrepo"
      mount_path  = "/gitrepo"
      read_only   = false
    }

    volume_mount {
      volume_name = "azureshare"
      mount_path  = "/azureshare"
      read_only   = false
    }
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

* `volume` - (Optional) The definition of volumes for use by containers via `volume_mounts` as documented in the `volume` block below.

The `volume` block supports:

* `name` - (Required) Sepcifies the name of the Volume. Changing this foces a new resource to be created. 

* `azure_share` - (Optional) The definition of a Azure Share based volume as documented in the `azure_share` block below. 

* `git_repo` - (Optional) The definition of a Git Repository based volume as documented in the `git_repo` block below. 

* `empty_dir` - (Optional) The definition of a Empty Dir based volume. 

* `secret` - (Optional) The definition of a Secret based volume as documented in the `secret` block blow.

~> **Note:** the `volume` block should contain one of either `azure_share`, `git_repo`, `empty_dir` or `secret`, if more than one is defined the resource will return an error.


The `secret` block supports:

* `name` - (Required) The name of the secret. This will be the filename the secret appears on disk as in the `mount_path`.

* `data` - (Required) The data containered in the secret in base64 encoding. 


The `git_repo` block supports:

* `repository` - (Required) The URL of the public git repository to mount. 

* `directory` - (Optional) The subdirectory to use in the git repository. 


The `azure_share` block supports:

* `share_name` - (Required) The Azure storage share that is to be mounted as a volume. This must be created on the storage account specified as above. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The Azure storage account from which the volume is to be mounted. Changing this forces a new resource to be created.

* `storage_account_key` - (Required) The access key for the Azure Storage account specified as above. Changing this forces a new resource to be created.


The `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `port` - (Optional) A public port for the container. Changing this forces a new resource to be created.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container. Changing this forces a new resource to be created.

~> **NOTE:** The field `command` has been deprecated in favor of `commands` to better match the API.

* `commands` - (Optional) A list of commands which should be run on the container. Changing this forces a new resource to be created.

* `volume_mount` - (Optional) The definition of a volume mount for this container as documented in the `volume_mount` block below. Changing this forces a new resource to be created.

The `volume_mount` block supports:

* `volume_name` - (Required) The name of the `volume` to mount into this container. Changing this forces a new resource to be created.

* `mount_path` - (Required) The path on which this volume is to be mounted. Changing this forces a new resource to be created.

* `read_only` - (Optional) Specify if the volume is to be mounted as read only or not. The default value is `false`. Changing this forces a new resource to be created.

The `image_registry_credential` block supports:

* `username` - (Required) The username with which to connect to the registry.

* `password` - (Required) The password with which to connect to the registry.

* `server` - (Required) The address to use to connect to the registry without protocol ("https"/"http"). For example: "myacr.acr.io" 

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
