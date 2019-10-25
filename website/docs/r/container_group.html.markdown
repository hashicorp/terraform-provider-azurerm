---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
sidebar_current: "docs-azurerm-resource-container-group"
description: |-
  Create as an Azure Container Group instance.
---

# azurerm_container_group

Manages as an Azure Container Group instance.

## Example Usage

This example provisions a Basic Container. Other examples of the `azurerm_container_group` resource can be found in [the `./examples/container-instance` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/container-instance).

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_group" "example" {
  name                = "example-continst"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ip_address_type     = "public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"

    ports {
      port     = 443
      protocol = "TCP"
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

* `identity` - (Optional) An `identity` block as defined below.

~> **Note:** managed identities are not supported for containers in virtual networks.

* `container` - (Required) The definition of a container that is part of the group as documented in the `container` block below. Changing this forces a new resource to be created.

* `os_type` - (Required) The OS for the container group. Allowed values are `Linux` and `Windows`. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported. Windows containers are not supported in virtual networks.

---

* `diagnostics` - (Optional) A `diagnostics` block as documented below.

* `dns_name_label` - (Optional) The DNS label/name for the container groups IP. Changing this forces a new resource to be created.

~> **Note:** DNS label/name is not supported when deploying to virtual networks.

* `ip_address_type` - (Optional) Specifies the ip address type of the container. `Public` or `Private`. Changing this forces a new resource to be created. If set to `Private`, `network_profile_id` also needs to be set.

~> **Note:** `dns_name_label`, `identity` and `os_type` set to `windows` are not compatible with `Private` `ip_address_type`

* `network_profile_id` - (Optional) Network profile ID for deploying to virtual network.

* `image_registry_credential` - (Optional) A `image_registry_credential` block as documented below. Changing this forces a new resource to be created.

* `restart_policy` - (Optional) Restart policy for the container group. Allowed values are `Always`, `Never`, `OnFailure`. Defaults to `Always`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

An `identity` block supports the following:

* `type` - (Required) The Managed Service Identity Type of this container group. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities. Changing this forces a new resource to be created.

~> **NOTE:** When `type` is set to `SystemAssigned`, identity the Principal ID can be retrieved after the container group has been created. See [documentation](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview) for more information.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`. Changing this forces a new resource to be created.

---

A `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `gpu` - (Optional) A `gpu` block as defined below. Changing this forces a new resource to be created.

~> **Note:** Gpu resources are currently only supported in Linux containers.

* `ports` - (Optional) A set of public ports for the container. Changing this forces a new resource to be created. Set as documented in the `ports` block below.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `secure_environment_variables` - (Optional) A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `readiness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `readiness_probe` block below. Changing this forces a new resource to be created.

* `liveness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `liveness_probe` block below. Changing this forces a new resource to be created.

* `command` - (Optional) A command line to be run on the container. Changing this forces a new resource to be created.

~> **NOTE:** The field `command` has been deprecated in favor of `commands` to better match the API.

* `commands` - (Optional) A list of commands which should be run on the container. Changing this forces a new resource to be created.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

---

A `diagnostics` block supports:

* `log_analytics` - (Required) A `log_analytics` block as defined below. Changing this forces a new resource to be created.

---

A `image_registry_credential` block supports:

* `username` - (Required) The username with which to connect to the registry. Changing this forces a new resource to be created.

* `password` - (Required) The password with which to connect to the registry. Changing this forces a new resource to be created.

* `server` - (Required) The address to use to connect to the registry without protocol ("https"/"http"). For example: "myacr.acr.io". Changing this forces a new resource to be created.

---

A `log_analytics` block supports:

* `log_type` - (Optional) The log type which should be used. Possible values are `ContainerInsights` and `ContainerInstanceLogs`. Changing this forces a new resource to be created.

* `workspace_id` - (Required) The Workspace ID of the Log Analytics Workspace. Changing this forces a new resource to be created.

* `workspace_key` - (Required) The Workspace Key of the Log Analytics Workspace. Changing this forces a new resource to be created.

* `metadata` - (Optional) Any metadata required for Log Analytics. Changing this forces a new resource to be created.

---

A `ports` block supports:

* `port` - (Required) The port number the container will expose. Changing this forces a new resource to be created.

* `protocol` - (Required) The network protocol associated with port. Possible values are `TCP` & `UDP`. Changing this forces a new resource to be created.

--

A `gpu` block supports:

* `count` - (Required) The number of GPUs which should be assigned to this container. Allowed values are `1`, `2`, or `4`. Changing this forces a new resource to be created.

* `sku` - (Required) The Sku which should be used for the GPU. Possible values are `K80`, `P100`, or `V100`. Changing this forces a new resource to be created.

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
