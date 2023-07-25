---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app"
description: |-
  Manages a Container App.
---

# azurerm_container_app

Manages a Container App.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "Example-Environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}
resource "azurerm_container_app" "example" {
  name                         = "example-app"
  container_app_environment_id = azurerm_container_app_environment.example.id
  resource_group_name          = azurerm_resource_group.example.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `container_app_environment_id` - (Required) The ID of the Container App Environment within which this Container App should exist. Changing this forces a new resource to be created.

* `name` - (Required) The name for this Container App. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Container App Environment is to be created. Changing this forces a new resource to be created.

* `revision_mode` - (Required) The revisions operational mode for the Container App. Possible values include `Single` and `Multiple`. In `Single` mode, a single revision is in operation at any given time. In `Multiple` mode, more than one revision can be active at a time and can be configured with load distribution via the `traffic_weight` block in the `ingress` configuration.

* `template` - (Required) A `template` block as detailed below.

---

* `dapr` - (Optional) A `dapr` block as detailed below.

* `identity` - (Optional) An `identity` block as detailed below.

* `ingress` - (Optional) An `ingress` block as detailed below.

* `registry` - (Optional) A `registry` block as detailed below.

* `secret` - (Optional) One or more `secret` block as detailed below.

* `tags` - (Optional) A mapping of tags to assign to the Container App.

---

A `secret` block supports the following:

* `name` - (Required) The Secret name.

* `value` - (Required) The value for this secret.

!> **Note:** Secrets cannot be removed from the service once added, attempting to do so will result in an error. Their values may be zeroed, i.e. set to `""`, but the named secret must persist. This is due to a technical limitation on the service which causes the service to become unmanageable. See [this issue](https://github.com/microsoft/azure-container-apps/issues/395) for more details.

---

A `template` block supports the following:

* `container` - (Required) One or more `container` blocks as detailed below.

* `max_replicas` - (Optional) The maximum number of replicas for this container.

* `min_replicas` - (Optional) The minimum number of replicas for this container.

* `revision_suffix` - (Optional) The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.

* `volume` - (Optional) A `volume` block as detailed below.

---

A `volume` block supports the following:

* `name` - (Required) The name of the volume.

* `storage_name` - (Optional) The name of the `AzureFile` storage.

* `storage_type` - (Optional) The type of storage volume. Possible values include `AzureFile` and `EmptyDir`. Defaults to `EmptyDir`.

---

A `container` block supports the following:

* `args` - (Optional) A list of extra arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - (Required) The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. 

~> **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`

* `env` - (Optional) One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App. 

~> **NOTE:** `ephemeral_storage` is currently in preview and not configurable at this time.

* `image` - (Required) The image to use to create the container.

* `liveness_probe` - (Optional) A `liveness_probe` block as detailed below.

* `memory` - (Required) The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1.0Gi`, `1.5Gi`, `2.0Gi`, `2.5Gi`, `3.0Gi`, `3.5Gi`, and `4.0Gi`. 

~> **NOTE:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`

* `name` - (Required) The name of the container

* `readiness_probe` - (Optional) A `readiness_probe` block as detailed below.

* `startup_probe` - (Optional) A `startup_probe` block as detailed below.

* `volume_mounts` - (Optional) A `volume_mounts` block as detailed below.

---

A `liveness_probe` block supports the following:

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The time in seconds to wait after the container has started before the probe is started.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are in the range `1` - `240`. Defaults to `10`.

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `termination_grace_period_seconds` -  The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

An `env` block supports the following:

* `name` - (Required) The name of the environment variable for the container.

* `secret_name` - (Optional) The name of the secret that contains the value for this environment variable.

* `value` - (Optional) The value for this environment variable.

~> **NOTE:** This value is ignored if `secret_name` is used

---

A `readiness_probe` block supports the following:

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - (Optional) The URI to use for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `success_count_threshold` - (Optional) The number of consecutive successful responses required to consider this probe as successful. Possible values are between `1` and `10`. Defaults to `3`.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

A `startup_probe` block supports the following:

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The value for the host header which should be sent with this probe. If unspecified, the IP Address of the Pod is used as the host header. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `termination_grace_period_seconds` -  The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

A `volume_mounts` block supports the following:

* `name` - (Required) The name of the Volume to be mounted in the container.

* `path` - (Required) The path in the container at which to mount this volume.

---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ingress` block supports the following:

* `allow_insecure_connections` - (Optional) Should this ingress allow insecure connections?

* `custom_domain` -  (Optional) One or more `custom_domain` block as detailed below.

* `fqdn` -  The FQDN of the ingress.

* `external_enabled` - (Optional) Is this an external Ingress.

* `target_port` - (Required) The target port on the container for the Ingress traffic.

* `traffic_weight` - (Required) A `traffic_weight` block as detailed below.

~> **Note:** `traffic_weight` can only be specified when `revision_mode` is set to `Multiple`.

* `transport` - (Optional) The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`. Defaults to `auto`

---

A `custom_domain` block supports the following:

* `certificate_binding_type` - (Optional) The Binding type. Possible values include `Disabled` and `SniEnabled`. Defaults to `Disabled`.

* `certificate_id` - (Required) The ID of the Container App Environment Certificate.

* `name` - (Required) The hostname of the Certificate. Must be the CN or a named SAN in the certificate.

---

A `traffic_weight` block supports the following:

~> **Note:** This block only applies when `revision_mode` is set to `Multiple`.

* `label` - (Optional) The label to apply to the revision as a name prefix for routing traffic.

* `latest_revision` - (Optional) This traffic Weight relates to the latest stable Container Revision.

* `revision_suffix` - (Optional) The suffix string to which this `traffic_weight` applies.

* `percentage` - (Required) The percentage of traffic which should be sent this revision.

~> **Note:** The cumulative values for `weight` must equal 100 exactly and explicitly, no default weights are assumed. 

---

A `dapr` block supports the following:

* `app_id` - (Required) The Dapr Application Identifier.

* `app_port` - (Optional) The port which the application is listening on. This is the same as the `ingress` port.

* `app_protocol` - (Optional) The protocol for the app. Possible values include `http` and `grpc`. Defaults to `http`.

---

A `registry` block supports the following:

* `server` - (Required) The hostname for the Container Registry.

The authentication details must also be supplied, `identity` and `username`/`password_secret_name` are mutually exclusive.

* `identity` - (Optional) Resource ID for the User Assigned Managed identity to use when pulling from the Container Registry.

* `password_secret_name` - (Optional) The name of the Secret Reference containing the password value for this user on the Container Registry, `username` must also be supplied.

* `username` - (Optional) The username to use for this Container Registry, `password_secret_name` must also be supplied..



## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App.

* `custom_domain_verification_id` - The ID of the Custom Domain Verification for this Container App.

* `latest_revision_fqdn` - The FQDN of the Latest Revision of the Container App.

* `latest_revision_name` - The name of the latest Container Revision.

* `location` - The location this Container App is deployed in. This is the same as the Environment in which it is deployed.

* `outbound_ip_addresses` - A list of the Public IP Addresses which the Container App uses for outbound network access.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App.
* `update` - (Defaults to 30 minutes) Used when updating the Container App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App.

## Import

A Container App can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/containerApps/myContainerApp"
```
