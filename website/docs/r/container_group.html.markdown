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
      port     = 80
      protocol = "TCP"
    }
    ports {
      port     = 443
      protocol = "TCP"
    }

    environment_variables = {
      "NODE_ENV" = "testing"
    }

    secure_environment_variables = {
      "ACCESS_KEY" = "secure_testing"
    }

    readiness_probe {
      exec = ["/bin/sh","-c","touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600"]
    }

    liveness_probe {
      exec = ["cat", "/tmp/healthy"]
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

  tags = {
    environment = "testing"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container Group. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block.

* `container` - (Required) The definition of a container that is part of the group as documented in the `container` block below. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported.

* `os_type` - (Required) The OS for the container group. Allowed values are `Linux` and `Windows`. Changing this forces a new resource to be created.

---

* `diagnostics` - (Optional) A `diagnostics` block as documented below.

* `dns_name_label` - (Optional) The DNS label/name for the container groups IP.

* `ip_address_type` - (Optional) Specifies the ip address type of the container. `Public` is the only acceptable value at this time. Changing this forces a new resource to be created.

* `image_registry_credential` - (Optional) A `image_registry_credential` block as documented below.

* `restart_policy` - (Optional) Restart policy for the container group. Allowed values are `Always`, `Never`, `OnFailure`. Defaults to `Always`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) The Managed Service Identity Type of this container group. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the container group has been created. See [documentation](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) for more information.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

A `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `gpu` - (Optional) A `gpu` block as defined below.

~> **Note:** Gpu resources are currently only supported in Linux containers.

* `ports` - (Optional) A set of public ports for the container. Changing this forces a new resource to be created. Set as documented in the `ports` block below.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `secure_environment_variables` - (Optional) A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `readiness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `readiness_probe` block below. Changing this forces a new resource to be created.

* `liveness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `liveness_probe` block below. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container.

~> **NOTE:** The field `command` has been deprecated in favor of `commands` to better match the API.

* `commands` - (Optional) A list of commands which should be run on the container.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

---

A `diagnostics` block supports:

* `log_analytics` - (Required) A `log_analytics` block as defined below.

---

A `image_registry_credential` block supports:

* `username` - (Required) The username with which to connect to the registry.

* `password` - (Required) The password with which to connect to the registry.

* `server` - (Required) The address to use to connect to the registry without protocol ("https"/"http"). For example: "myacr.acr.io"

---

A `log_analytics` block supports:

* `log_type` - (Required) The log type which should be used. Possible values are `ContainerInsights` and `ContainerInstanceLogs`.

* `workspace_id` - (Required) The Workspace ID of the Log Analytics Workspace.

* `workspace_key` - (Required) The Workspace Key of the Log Analytics Workspace.

* `metadata` - (Optional) Any metadata required for Log Analytics.

---

A `ports` block supports:

* `port` - (Required) The port number the container will expose.

* `protocol` - (Required) The network protocol associated with port. Possible values are `TCP` & `UDP`.

--

A `gpu` block supports:

* `count` - (Required) The number of GPUs which should be assigned to this container. Allowed values are `1`, `2`, or `4`.

* `sku` - (Required) The Sku which should be used for the GPU. Possible values are `K80`, `P100`, or `V100`.

---

A `volume` block supports:

* `name` - (Required) The name of the volume mount. Changing this forces a new resource to be created.

* `mount_path` - (Required) The path on which this volume is to be mounted. Changing this forces a new resource to be created.

* `read_only` - (Optional) Specify if the volume is to be mounted as read only or not. The default value is `false`. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The Azure storage account from which the volume is to be mounted. Changing this forces a new resource to be created.

* `storage_account_key` - (Required) The access key for the Azure Storage account specified as above. Changing this forces a new resource to be created.

* `share_name` - (Required) The Azure storage share that is to be mounted as a volume. This must be created on the storage account specified as above. Changing this forces a new resource to be created.

---

The `readiness_probe` block supports:

* `exec` - (Optional) Commands to be run to validate container readiness. Changing this forces a new resource to be created.

* `httpget` - (Optional) The definition of the httpget for this container as documented in the `httpget` block below. Changing this forces a new resource to be created.

* `initial_delay_seconds` - (Optional) Number of seconds after the container has started before liveness or readiness probes are initiated. Changing this forces a new resource to be created.

* `period_seconds` - (Optional) How often (in seconds) to perform the probe. The default value is `10` and the minimum value is `1`. Changing this forces a new resource to be created.

* `failure_threshold` - (Optional) How many times to try the probe before restarting the container (liveness probe) or marking the container as unhealthy (readiness probe). The default value is `3` and the minimum value is `1`. Changing this forces a new resource to be created.

* `success_threshold` - (Optional) Minimum consecutive successes for the probe to be considered successful after having failed. The default value is `1` and the minimum value is `1`. Changing this forces a new resource to be created.

* `timeout_seconds` - (Optional) Number of seconds after which the probe times out. The default value is `1` and the minimum value is `1`. Changing this forces a new resource to be created.

---

The `liveness_probe` block supports:

* `exec` - (Optional) Commands to be run to validate container readiness. Changing this forces a new resource to be created.

* `httpget` - (Optional) The definition of the httpget for this container as documented in the `httpget` block below. Changing this forces a new resource to be created.

* `initial_delay_seconds` - (Optional) Number of seconds after the container has started before liveness or readiness probes are initiated. Changing this forces a new resource to be created.

* `period_seconds` - (Optional) How often (in seconds) to perform the probe. The default value is `10` and the minimum value is `1`. Changing this forces a new resource to be created.

* `failure_threshold` - (Optional) How many times to try the probe before restarting the container (liveness probe) or marking the container as unhealthy (readiness probe). The default value is `3` and the minimum value is `1`. Changing this forces a new resource to be created.

* `success_threshold` - (Optional) Minimum consecutive successes for the probe to be considered successful after having failed. The default value is `1` and the minimum value is `1`. Changing this forces a new resource to be created.

* `timeout_seconds` - (Optional) Number of seconds after which the probe times out. The default value is `1` and the minimum value is `1`. Changing this forces a new resource to be created.

---

The `httpget` block supports:

* `path` - (Optional) Path to access on the HTTP server. Changing this forces a new resource to be created.

* `port` - (Optional) Number of the port to access on the container. Changing this forces a new resource to be created.

* `scheme` - (Optional) Scheme to use for connecting to the host. Possible values are `Http` and `Https`. Changing this forces a new resource to be created.

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
