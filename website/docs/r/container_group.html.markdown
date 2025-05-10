---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_group"
description: |-
  Manages an Azure Container Group instance.
---

# azurerm_container_group

Manages as an Azure Container Group instance.

## Example Usage

This example provisions a Basic Container. Other examples of the `azurerm_container_group` resource can be found in [the `./examples/container-instance` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/container-instance).

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_container_group" "example" {
  name                = "example-continst"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"

  container {
    name   = "hello-world"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"

    ports {
      port     = 443
      protocol = "TCP"
    }
  }

  container {
    name   = "sidecar"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
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

* `sku` - (Optional) Specifies the sku of the Container Group. Possible values are `Confidential`, `Dedicated` and `Standard`. Defaults to `Standard`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `init_container` - (Optional) The definition of an init container that is part of the group as documented in the `init_container` block below. Changing this forces a new resource to be created.

* `container` - (Required) The definition of a container that is part of the group as documented in the `container` block below. Changing this forces a new resource to be created.

* `os_type` - (Required) The OS for the container group. Allowed values are `Linux` and `Windows`. Changing this forces a new resource to be created.

~> **Note:** if `os_type` is set to `Windows` currently only a single `container` block is supported. Windows containers are not supported in virtual networks.

---

* `dns_config` - (Optional) A `dns_config` block as documented below. Changing this forces a new resource to be created.

* `diagnostics` - (Optional) A `diagnostics` block as documented below. Changing this forces a new resource to be created.

* `dns_name_label` - (Optional) The DNS label/name for the container group's IP. Changing this forces a new resource to be created.

~> **Note:** DNS label/name is not supported when deploying to virtual networks.

* `dns_name_label_reuse_policy` - (Optional) The value representing the security enum. `Noreuse`, `ResourceGroupReuse`, `SubscriptionReuse`, `TenantReuse` or `Unsecure`. Defaults to `Unsecure`. Changing this forces a new resource to be created.

* `exposed_port` - (Optional) Zero or more `exposed_port` blocks as defined below. Changing this forces a new resource to be created.

~> **Note:** The `exposed_port` can only contain ports that are also exposed on one or more containers in the group.

* `ip_address_type` - (Optional) Specifies the IP address type of the container. `Public`, `Private` or `None`. Changing this forces a new resource to be created. If set to `Private`, `subnet_ids` also needs to be set. Defaults to `Public`.

~> **Note:** `dns_name_label` and `os_type` set to `windows` are not compatible with `Private` `ip_address_type`

* `key_vault_key_id` - (Optional) The Key Vault key URI for CMK encryption. Changing this forces a new resource to be created.

* `key_vault_user_assigned_identity_id` - (Optional) The user assigned identity that has access to the Key Vault Key. If not specified, the RP principal named "Azure Container Instance Service" will be used instead. Make sure the identity has the proper `key_permissions` set, at least with `Get`, `UnwrapKey`, `WrapKey` and `GetRotationPolicy`.

* `subnet_ids` - (Optional) The subnet resource IDs for a container group. Changing this forces a new resource to be created.

* `image_registry_credential` - (Optional) An `image_registry_credential` block as documented below. Changing this forces a new resource to be created.

* `priority` - (Optional) The priority of the Container Group. Possible values are `Regular` and `Spot`. Changing this forces a new resource to be created.

~> **Note:** When `priority` is set to `Spot`, the `ip_address_type` has to be `None`.

* `restart_policy` - (Optional) Restart policy for the container group. Allowed values are `Always`, `Never`, `OnFailure`. Defaults to `Always`. Changing this forces a new resource to be created.

* `zones` - (Optional) A list of Availability Zones in which this Container Group is located. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Container Group. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

~> **Note:** When `type` is set to `SystemAssigned`, the identity of the Principal ID can be retrieved after the container group has been created. See [documentation](https://docs.microsoft.com/azure/active-directory/managed-service-identity/overview) for more information.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Container Group.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

~> **Note:** Currently you can't use a managed identity in a container group deployed to a virtual network.

---

An `init_container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `secure_environment_variables` - (Optional) A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `commands` - (Optional) A list of commands which should be run on the container. Changing this forces a new resource to be created.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

* `security` - (Optional) The definition of the security context for this container as documented in the `security` block below. Changing this forces a new resource to be created.

---

A `container` block supports:

* `name` - (Required) Specifies the name of the Container. Changing this forces a new resource to be created.

* `image` - (Required) The container image name. Changing this forces a new resource to be created.

* `cpu` - (Required) The required number of CPU cores of the containers. Changing this forces a new resource to be created.

* `memory` - (Required) The required memory of the containers in GB. Changing this forces a new resource to be created.

* `cpu_limit` - (Optional) The upper limit of the number of CPU cores of the containers.

* `memory_limit` - (Optional) The upper limit of the memory of the containers in GB.

* `ports` - (Optional) A set of public ports for the container. Changing this forces a new resource to be created. Set as documented in the `ports` block below.

* `environment_variables` - (Optional) A list of environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `secure_environment_variables` - (Optional) A list of sensitive environment variables to be set on the container. Specified as a map of name/value pairs. Changing this forces a new resource to be created.

* `readiness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `readiness_probe` block below. Changing this forces a new resource to be created.

* `liveness_probe` - (Optional) The definition of a readiness probe for this container as documented in the `liveness_probe` block below. Changing this forces a new resource to be created.

* `commands` - (Optional) A list of commands which should be run on the container. Changing this forces a new resource to be created.

* `volume` - (Optional) The definition of a volume mount for this container as documented in the `volume` block below. Changing this forces a new resource to be created.

* `security` - (Optional) The definition of the security context for this container as documented in the `security` block below. Changing this forces a new resource to be created.

---

An `exposed_port` block supports:

* `port` - (Optional) The port number the container will expose. Changing this forces a new resource to be created.

* `protocol` - (Optional) The network protocol associated with port. Possible values are `TCP` & `UDP`. Changing this forces a new resource to be created. Defaults to `TCP`.

~> **Note:** Removing all `exposed_port` blocks requires setting `exposed_port = []`.

---

A `diagnostics` block supports:

* `log_analytics` - (Required) A `log_analytics` block as defined below. Changing this forces a new resource to be created.

---

An `image_registry_credential` block supports:

* `user_assigned_identity_id` - (Optional) The identity ID for the private registry. Changing this forces a new resource to be created.

* `username` - (Optional) The username with which to connect to the registry. Changing this forces a new resource to be created.

* `password` - (Optional) The password with which to connect to the registry. Changing this forces a new resource to be created.

* `server` - (Required) The address to use to connect to the registry without protocol ("https"/"http"). For example: "myacr.acr.io". Changing this forces a new resource to be created.

---

A `log_analytics` block supports:

* `log_type` - (Optional) The log type which should be used. Possible values are `ContainerInsights` and `ContainerInstanceLogs`. Changing this forces a new resource to be created.

* `workspace_id` - (Required) The Workspace ID of the Log Analytics Workspace. Changing this forces a new resource to be created.

* `workspace_key` - (Required) The Workspace Key of the Log Analytics Workspace. Changing this forces a new resource to be created.

* `metadata` - (Optional) Any metadata required for Log Analytics. Changing this forces a new resource to be created.

---

A `ports` block supports:

* `port` - (Optional) The port number the container will expose. Changing this forces a new resource to be created.

* `protocol` - (Optional) The network protocol associated with port. Possible values are `TCP` & `UDP`. Changing this forces a new resource to be created. Defaults to `TCP`.

~> **Note:** Omitting these blocks will default the exposed ports on the group to all ports on all containers defined in the `container` blocks of this group.

---

A `volume` block supports:

* `name` - (Required) The name of the volume mount. Changing this forces a new resource to be created.

* `mount_path` - (Required) The path on which this volume is to be mounted. Changing this forces a new resource to be created.

* `read_only` - (Optional) Specify if the volume is to be mounted as read only or not. The default value is `false`. Changing this forces a new resource to be created.

* `empty_dir` - (Optional) Boolean as to whether the mounted volume should be an empty directory. Defaults to `false`. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The Azure storage account from which the volume is to be mounted. Changing this forces a new resource to be created.

* `storage_account_key` - (Optional) The access key for the Azure Storage account specified as above. Changing this forces a new resource to be created.

* `share_name` - (Optional) The Azure storage share that is to be mounted as a volume. This must be created on the storage account specified as above. Changing this forces a new resource to be created.

* `git_repo` - (Optional) A `git_repo` block as defined below. Changing this forces a new resource to be created.

* `secret` - (Optional) A map of secrets that will be mounted as files in the volume. Changing this forces a new resource to be created.

~> **Note:** Exactly one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) must be specified.

~> **Note:** when using a storage account volume, all of `share_name`, `storage_account_name`, and `storage_account_key` must be specified.

~> **Note:** The secret values must be supplied as Base64 encoded strings, such as by using the Terraform [base64encode function](https://www.terraform.io/docs/configuration/functions/base64encode.html). The secret values are decoded to their original values when mounted in the volume on the container.

---

The `git_repo` block supports:

* `url` - (Required) Specifies the Git repository to be cloned. Changing this forces a new resource to be created.

* `directory` - (Optional) Specifies the directory into which the repository should be cloned. Changing this forces a new resource to be created.

* `revision` - (Optional) Specifies the commit hash of the revision to be cloned. If unspecified, the HEAD revision is cloned. Changing this forces a new resource to be created.

---

The `readiness_probe` block supports:

* `exec` - (Optional) Commands to be run to validate container readiness. Changing this forces a new resource to be created.

* `http_get` - (Optional) The definition of the http_get for this container as documented in the `http_get` block below. Changing this forces a new resource to be created.

* `initial_delay_seconds` - (Optional) Number of seconds after the container has started before liveness or readiness probes are initiated. Changing this forces a new resource to be created.

* `period_seconds` - (Optional) How often (in seconds) to perform the probe. Changing this forces a new resource to be created.

* `failure_threshold` - (Optional) How many times to try the probe before restarting the container (liveness probe) or marking the container as unhealthy (readiness probe). Changing this forces a new resource to be created.

* `success_threshold` - (Optional) Minimum consecutive successes for the probe to be considered successful after having failed. Changing this forces a new resource to be created.

* `timeout_seconds` - (Optional) Number of seconds after which the probe times out. Changing this forces a new resource to be created.

---

The `liveness_probe` block supports:

* `exec` - (Optional) Commands to be run to validate container readiness. Changing this forces a new resource to be created.

* `http_get` - (Optional) The definition of the http_get for this container as documented in the `http_get` block below. Changing this forces a new resource to be created.

* `initial_delay_seconds` - (Optional) Number of seconds after the container has started before liveness or readiness probes are initiated. Changing this forces a new resource to be created.

* `period_seconds` - (Optional) How often (in seconds) to perform the probe. Changing this forces a new resource to be created.

* `failure_threshold` - (Optional) How many times to try the probe before restarting the container (liveness probe) or marking the container as unhealthy (readiness probe). Changing this forces a new resource to be created.

* `success_threshold` - (Optional) Minimum consecutive successes for the probe to be considered successful after having failed. Changing this forces a new resource to be created.

* `timeout_seconds` - (Optional) Number of seconds after which the probe times out. Changing this forces a new resource to be created.

---

The `http_get` block supports:

* `path` - (Optional) Path to access on the HTTP server. Changing this forces a new resource to be created.

* `port` - (Optional) Number of the port to access on the container. Changing this forces a new resource to be created.

* `scheme` - (Optional) Scheme to use for connecting to the host. Possible values are `Http` and `Https`. Changing this forces a new resource to be created.

* `http_headers` - (Optional) A map of HTTP headers used to access on the container. Changing this forces a new resource to be created.

---

The `dns_config` block supports:

* `nameservers` - (Required) A list of nameservers the containers will search out to resolve requests. Changing this forces a new resource to be created.

* `search_domains` - (Optional) A list of search domains that DNS requests will search along. Changing this forces a new resource to be created.

* `options` - (Optional) A list of [resolver configuration options](https://man7.org/linux/man-pages/man5/resolv.conf.5.html). Changing this forces a new resource to be created.

---

The `security` block supports:

* `privilege_enabled` - (Required) Whether the container's permission is elevated to privileged? Changing this forces a new resource to be created.

~> **Note:** Currently, this only applies when the `os_type` is `Linux` and the `sku` is `Confidential`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Group.

* `identity` - An `identity` block as defined below.

* `ip_address` - The IP address allocated to the container group.

* `fqdn` - The FQDN of the container group derived from `dns_name_label`.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Container Group.

* `read` - (Defaults to 5 minutes) Used when retrieving the Container Group.

* `update` - (Defaults to 30 minutes) Used when updating the Container Group.

* `delete` - (Defaults to 30 minutes) Used when deleting the Container Group.

## Import

Container Group's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_group.containerGroup1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ContainerInstance/containerGroups/myContainerGroup1
```
