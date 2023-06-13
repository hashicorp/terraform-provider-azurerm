---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app"
description: |-
  Get information of a Container App.
---

# Data Source: azurerm_container_app

Use this data source to access information about an existing Container App.

## Example Usage

```hcl
data "azurerm_container_app" "example" {
  name                = "example-app"
  resource_group_name = "example-resources"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Container App.

* `resource_group_name` - (Required) The name of the Resource Group where this Container App exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `container_app_environment_id` - The ID of the Container App Environment this Container App is linked to.

* `revision_mode` - The revision mode of the Container App.

* `template` - A `template` block as detailed below.

---

* `dapr` - A `dapr` block as detailed below.

* `identity` - An `identity` block as detailed below.

* `ingress` - An `ingress` block as detailed below.

* `registry` - A `registry` block as detailed below.

* `secret` - One or more `secret` block as detailed below.

* `tags` - A mapping of tags to assign to the Container App.

---

A `secret` block supports the following:

* `name` - The Secret name.

* `value` - The value for this secret.

---

A `template` block supports the following:

* `container` - One or more `container` blocks as detailed below.

* `max_replicas` - The maximum number of replicas for this container.

* `min_replicas` - The minimum number of replicas for this container.

* `revision_suffix` - The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.

* `volume` - A `volume` block as detailed below.

---

A `volume` block supports the following:

* `name` - The name of the volume.

* `storage_name` - The name of the `AzureFile` storage.

* `storage_type` - The type of storage volume. Possible values include `AzureFile` and `EmptyDir`. Defaults to `EmptyDir`.

---

A `container` block supports the following:

* `args` - A list of extra arguments to pass to the container.

* `command` - A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`.

* `env` - One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

* `image` - The image to use to create the container.

* `liveness_probe` - A `liveness_probe` block as detailed below.

* `memory` - The amount of memory to allocate to the container. Possible values include `0.5Gi`, `1Gi`, `1.5Gi`, `2Gi`, `2.5Gi`, `3Gi`, `3.5Gi`, and `4Gi`.

* `name` - The name of the container

* `readiness_probe` - A `readiness_probe` block as detailed below.

* `startup_probe` - A `startup_probe` block as detailed below.

* `volume_mounts` - A `volume_mounts` block as detailed below.

---

A `liveness_probe` block supports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - A `header` block as detailed below.

* `host` - The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - The time in seconds to wait after the container has started before the probe is started.

* `interval_seconds` - How often, in seconds, the probe should run. Possible values are in the range `1` - `240`. Defaults to `10`.

* `path` - The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - The port number on which to connect. Possible values are between `1` and `65535`.

* `termination_grace_period_seconds` -  The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `timeout` - Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

An `env` block supports the following:

* `name` - The name of the environment variable for the container.

* `secret_name` - The name of the secret that contains the value for this environment variable.

* `value` - The value for this environment variable.

---

A `readiness_probe` block supports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - A `header` block as detailed below.

* `host` - The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `interval_seconds` - How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - The URI to use for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - The port number on which to connect. Possible values are between `1` and `65535`.

* `success_count_threshold` - The number of consecutive successful responses required to consider this probe as successful. Possible values are between `1` and `10`. Defaults to `3`.

* `timeout` - Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

A `startup_probe` block supports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - A `header` block as detailed below.

* `host` - The value for the host header which should be sent with this probe. If unspecified, the IP Address of the Pod is used as the host header. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `interval_seconds` - How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - The port number on which to connect. Possible values are between `1` and `65535`.

* `termination_grace_period_seconds` -  The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `timeout` - Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

* `transport` - Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

---

A `header` block supports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

A `volume_mounts` block supports the following:

* `name` - The name of the Volume to be mounted in the container.

* `path` - The path in the container at which to mount this volume.

---

An `identity` block supports the following:

* `type` - The type of managed identity to assign. Possible values are `UserAssigned` and `SystemAssigned`

* `identity_ids` - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned`.

---

An `ingress` block supports the following:

* `allow_insecure_connections` - Should this ingress allow insecure connections?

* `custom_domain` - One or more `custom_domain` block as detailed below.

* `fqdn` -  The FQDN of the ingress.

* `external_enabled` - Is this an external Ingress.

* `target_port` - The target port on the container for the Ingress traffic.

* `traffic_weight` - A `traffic_weight` block as detailed below.

* `transport` - The transport method for the Ingress. Possible values include `auto`, `http`, and `http2`. Defaults to `auto`

---

A `custom_domain` block supports the following:

* `certificate_binding_type` - The Binding type. Possible values include `Disabled` and `SniEnabled`. Defaults to `Disabled`.

* `certificate_id` - The ID of the Container App Environment Certificate.

* `name` - The hostname of the Certificate. Must be the CN or a named SAN in the certificate.

---

A `traffic_weight` block supports the following:

* `label` - The label to apply to the revision as a name prefix for routing traffic.

* `latest_revision` - This traffic Weight relates to the latest stable Container Revision.

* `revision_suffix` - The suffix string to which this `traffic_weight` applies.

* `percentage` - The percentage of traffic which should be sent this revision.

---

A `dapr` block supports the following:

* `app_id` - The Dapr Application Identifier.

* `app_port` - The port which the application is listening on. This is the same as the `ingress` port.

* `app_protocol` - The protocol for the app. Possible values include `http` and `grpc`. Defaults to `http`.

---

A `registry` block supports the following:

* `server` - The hostname for the Container Registry.

* `identity` - Resource ID for the User Assigned Managed identity to use when pulling from the Container Registry.

* `password_secret_name` - The name of the Secret Reference containing the password value for this user on the Container Registry, `username` must also be supplied.

* `username` - The username to use for this Container Registry, `password_secret_name` must also be supplied..

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App.
