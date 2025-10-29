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

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container App.

* `resource_group_name` - (Required) The name of the Resource Group where this Container App exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `container_app_environment_id` - The ID of the Container App Environment this Container App is linked to.

* `revision_mode` - The revision mode of the Container App.

* `template` - A `template` block as detailed below.

* `dapr` - A `dapr` block as detailed below.

* `identity` - An `identity` block as detailed below.

* `ingress` - An `ingress` block as detailed below.

* `registry` - A `registry` block as detailed below.

* `secret` - One or more `secret` block as detailed below.

* `workload_profile_name` - The name of the Workload Profile in the Container App Environment in which this Container App is running.

* `max_inactive_revisions` - The max inactive revisions for this Container App.

* `tags` - A mapping of tags to assign to the Container App.

---

A `secret` block exports the following:

* `identity` - The identity used for accessing the Key Vault.

* `key_vault_secret_id` - The ID of a Key Vault secret.

* `name` - The secret name.

* `value` - The value for this secret.

---

A `template` block exports the following:

* `init_container` - One or more `init_container` blocks as detailed below.

* `container` - One or more `container` blocks as detailed below.

* `max_replicas` - The maximum number of replicas for this container.

* `min_replicas` - The minimum number of replicas for this container.

* `revision_suffix` - The suffix for the revision.

* `termination_grace_period_seconds` - The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `volume` - A `volume` block as detailed below.

---

A `volume` block exports the following:

* `name` - The name of the volume.

* `storage_name` - The name of the `AzureFile` storage.

* `storage_type` - The type of storage volume.

* `mount_options` - Mount options used while mounting the AzureFile.

* `secret` - A `secret` block as detailed below.

---

A `secret` block supports the following:

* `path` - Mount path for the secret

* `name` - Reference by name to the secret

---

A `init_container` block exports the following:

* `args` - A list of extra arguments passed to the container.

* `command` - A command passed to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - The amount of vCPU allocated to the container.

* `env` - One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

* `image` - The image to use to create the container.

* `memory` - The amount of memory allocated to the container.

* `name` - The name of the container

* `volume_mounts` - A `volume_mounts` block as detailed below.

---

A `container` block exports the following:

* `args` - A list of extra arguments passed to the container.

* `command` - A command passed to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - The amount of vCPU allocated to the container.

* `env` - One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

* `image` - The image to use to create the container.

* `liveness_probe` - A `liveness_probe` block as detailed below.

* `memory` - The amount of memory allocated to the container.

* `name` - The name of the container

* `readiness_probe` - A `readiness_probe` block as detailed below.

* `startup_probe` - A `startup_probe` block as detailed below.

* `volume_mounts` - A `volume_mounts` block as detailed below.

---

A `liveness_probe` block exports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed.

* `header` - A `header` block as detailed below.

* `host` - The probe hostname.

* `initial_delay` - The number of seconds elapsed after the container has started before the probe is initiated.

* `interval_seconds` - How often, in seconds, the probe should run.

* `path` - The URI used with the `host` for http type probes.

* `port` - The port number on which to connect.

* `timeout` - Time in seconds after which the probe times out.

* `transport` - Type of probe.

---

A `header` block exports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

An `env` block exports the following:

* `name` - The name of the environment variable for the container.

* `secret_name` - The name of the secret that contains the value for this environment variable.

* `value` - The value for this environment variable.

---

A `readiness_probe` block exports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed.

* `header` - A `header` block as detailed below.

* `host` - The probe hostname.

* `initial_delay` - The number of seconds elapsed after the container has started before the probe is initiated.

* `interval_seconds` - How often, in seconds, the probe should run.

* `path` - The URI to use for http type probes.

* `port` - The port number on which to connect.

* `success_count_threshold` - The number of consecutive successful responses required to consider this probe as successful.

* `timeout` - Time in seconds after which the probe times out.

* `transport` - Type of probe.

---

A `header` block exports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

A `startup_probe` block exports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed.

* `header` - A `header` block as detailed below.

* `host` - The value for the host header which should be sent with this probe.

* `initial_delay` - The number of seconds elapsed after the container has started before the probe is initiated.

* `interval_seconds` - How often, in seconds, the probe should run.

* `path` - The URI to use with the `host` for http type probes.

* `port` - The port number on which to connect.

* `timeout` - Time in seconds after which the probe times out.

* `transport` - Type of probe.

---

A `header` block exports the following:

* `name` - The HTTP Header Name.

* `value` - The HTTP Header value.

---

A `volume_mounts` block exports the following:

* `name` - The name of the Volume to be mounted in the container.

* `path` - The path in the container at which to mount this volume.

* `sub_path` - The sub path of the volume to be mounted in the container.

---

An `identity` block exports the following:

* `type` - The type of managed identity to assign.

* `identity_ids` - A list of one or more Resource IDs for User Assigned Managed identities to assign.

---

An `ingress` block exports the following:

* `allow_insecure_connections` - Should this ingress allow insecure connections?

* `client_certificate_mode` - The client certificate mode for the Ingress.

* `cors` - A `cors` block as detailed below.

* `custom_domain` - One or more `custom_domain` block as detailed below.

* `fqdn` -  The FQDN of the ingress.

* `external_enabled` - Is this an external Ingress.

* `ip_security_restriction` - One or more `ip_security_restriction` blocks for IP-filtering rules as defined below.

* `target_port` - The target port on the container for the Ingress traffic.

* `traffic_weight` - A `traffic_weight` block as detailed below.

* `transport` - The transport method for the Ingress.

---

A `cors` block exports the following:

* `allowed_origins` - The list of origins that are allowed to make cross-origin calls.

* `allow_credentials_enabled` - Whether user credentials are allowed in the cross-origin request.

* `allowed_headers` - The list of request headers that are permitted in the actual request.

* `allowed_methods` - The list of HTTP methods are allowed when accessing the resource in a cross-origin request.

* `exposed_headers` - The list of headers exposed to the browser in the response to a cross-origin request.

* `max_age_in_seconds` - The number of seconds that the browser can cache the results of a preflight request.

---

A `custom_domain` block exports the following:

* `certificate_binding_type` - The Binding type.

* `certificate_id` - The ID of the Container App Environment Certificate.

* `name` - The hostname of the Certificate. Must be the CN or a named SAN in the certificate.

---

A `ip_security_restriction` block exports the following:

* `action` - The IP-filter action.

* `description` - Description of the IP restriction rule that is being sent to the container-app.

* `ip_address_range` - CIDR notation that matches the incoming IP address.

* `name` - Name for the IP restriction rule.

---

A `traffic_weight` block exports the following:

* `label` - The label to apply to the revision as a name prefix for routing traffic.

* `latest_revision` - This traffic Weight relates to the latest stable Container Revision.

* `revision_suffix` - The suffix string to which this `traffic_weight` applies.

* `percentage` - The percentage of traffic which should be sent this revision.

---

A `dapr` block exports the following:

* `app_id` - The Dapr Application Identifier.

* `app_port` - The port which the application is listening on. This is the same as the `ingress` port.

* `app_protocol` - The protocol for the app.

---

A `registry` block exports the following:

* `server` - The hostname for the Container Registry.

* `identity` - Resource ID for the User Assigned Managed identity to use when pulling from the Container Registry.

* `password_secret_name` - The name of the Secret Reference containing the password value for the user on the Container Registry.

* `username` - The username used for this Container Registry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
