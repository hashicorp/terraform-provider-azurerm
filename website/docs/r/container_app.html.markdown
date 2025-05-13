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
      image  = "mcr.microsoft.com/k8se/quickstart:latest"
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

* `workload_profile_name` - (Optional) The name of the Workload Profile in the Container App Environment to place this Container App.

~> **Note:** Omit this value to use the default `Consumption` Workload Profile.

* `max_inactive_revisions` - (Optional) The maximum of inactive revisions allowed for this Container App.

* `tags` - (Optional) A mapping of tags to assign to the Container App.

---

A `secret` block supports the following:

* `name` - (Required) The secret name.

* `identity` - (Optional) The identity to use for accessing the Key Vault secret reference. This can either be the Resource ID of a User Assigned Identity, or `System` for the System Assigned Identity.

!> **Note:** `identity` must be used together with `key_vault_secret_id`

* `key_vault_secret_id` - (Optional) The ID of a Key Vault secret. This can be a versioned or version-less ID.

!> **Note:** When using `key_vault_secret_id`, `ignore_changes` should be used to ignore any changes to `value`.

* `value` - (Optional) The value for this secret.

!> **Note:** `value` will be ignored if `key_vault_secret_id` and `identity` are provided.

---

A `template` block supports the following:

* `init_container` - (Optional) The definition of an init container that is part of the group as documented in the `init_container` block below.

* `container` - (Required) One or more `container` blocks as detailed below.

* `max_replicas` - (Optional) The maximum number of replicas for this container.

* `min_replicas` - (Optional) The minimum number of replicas for this container.

* `azure_queue_scale_rule` - (Optional) One or more `azure_queue_scale_rule` blocks as defined below.

* `custom_scale_rule` - (Optional) One or more `custom_scale_rule` blocks as defined below.

* `http_scale_rule` - (Optional) One or more `http_scale_rule` blocks as defined below.

* `tcp_scale_rule` - (Optional) One or more `tcp_scale_rule` blocks as defined below.

* `revision_suffix` - (Optional) The suffix for the revision. This value must be unique for the lifetime of the Resource. If omitted the service will use a hash function to create one.

* `termination_grace_period_seconds` - (Optional)   The time in seconds after the container is sent the termination signal before the process if forcibly killed.

* `volume` - (Optional) A `volume` block as detailed below.

---

An `azure_queue_scale_rule` block supports the following:

* `name` - (Required) The name of the Scaling Rule

* `queue_name` - (Required) The name of the Azure Queue

* `queue_length` - (Required) The value of the length of the queue to trigger scaling actions.

* `authentication` - (Required) One or more `authentication` blocks as defined below.

---

A `custom_scale_rule` block supports the following:

* `name` - (Required) The name of the Scaling Rule

* `custom_rule_type` - (Required) The Custom rule type. Possible values include: `activemq`, `artemis-queue`, `kafka`, `pulsar`, `aws-cloudwatch`, `aws-dynamodb`, `aws-dynamodb-streams`, `aws-kinesis-stream`, `aws-sqs-queue`, `azure-app-insights`, `azure-blob`, `azure-data-explorer`, `azure-eventhub`, `azure-log-analytics`, `azure-monitor`, `azure-pipelines`, `azure-servicebus`, `azure-queue`, `cassandra`, `cpu`, `cron`, `datadog`, `elasticsearch`, `external`, `external-push`, `gcp-stackdriver`, `gcp-storage`, `gcp-pubsub`, `graphite`, `http`, `huawei-cloudeye`, `ibmmq`, `influxdb`, `kubernetes-workload`, `liiklus`, `memory`, `metrics-api`, `mongodb`, `mssql`, `mysql`, `nats-jetstream`, `stan`, `tcp`, `new-relic`, `openstack-metric`, `openstack-swift`, `postgresql`, `predictkube`, `prometheus`, `rabbitmq`, `redis`, `redis-cluster`, `redis-sentinel`, `redis-streams`, `redis-cluster-streams`, `redis-sentinel-streams`, `selenium-grid`,`solace-event-queue`, and `github-runner`.

* `metadata` - (Required) - A map of string key-value pairs to configure the Custom Scale Rule.

* `authentication` - (Optional) Zero or more `authentication` blocks as defined below.

---

A `http_scale_rule` block supports the following:

* `name` - (Required) The name of the Scaling Rule

* `concurrent_requests` - (Required) - The number of concurrent requests to trigger scaling.

* `authentication` - (Optional) Zero or more `authentication` blocks as defined below.

---

A `tcp_scale_rule` block supports the following:

* `name` - (Required) The name of the Scaling Rule

* `concurrent_requests` - (Required) - The number of concurrent requests to trigger scaling.

* `authentication` - (Optional) Zero or more `authentication` blocks as defined below.

---

An `authentication` block supports the following:

* `secret_name` - (Required) The name of the Container App Secret to use for this Scale Rule Authentication.

* `trigger_parameter` - (Required) The Trigger Parameter name to use the supply the value retrieved from the `secret_name`.

---

A `volume` block supports the following:

* `name` - (Required) The name of the volume.

* `storage_name` - (Optional) The name of the `AzureFile` storage.

* `storage_type` - (Optional) The type of storage volume. Possible values are `AzureFile`, `EmptyDir` and `Secret`. Defaults to `EmptyDir`.

* `mount_options` - Mount options used while mounting the AzureFile. Must be a comma-separated string e.g. `dir_mode=0751,file_mode=0751`.

---

An `init_container` block supports:

* `args` - (Optional) A list of extra arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - (Optional) The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. When there's a workload profile specified, there's no such constraint.

~> **Note:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`

* `env` - (Optional) One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

~> **Note:** `ephemeral_storage` is currently in preview and not configurable at this time.

* `image` - (Required) The image to use to create the container.

* `memory` - (Optional) The amount of memory to allocate to the container. Possible values are `0.5Gi`, `1Gi`, `1.5Gi`, `2Gi`, `2.5Gi`, `3Gi`, `3.5Gi` and `4Gi`. When there's a workload profile specified, there's no such constraint.

~> **Note:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`

* `name` - (Required) The name of the container

* `volume_mounts` - (Optional) A `volume_mounts` block as detailed below.

---

A `container` block supports the following:

* `args` - (Optional) A list of extra arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `cpu` - (Required) The amount of vCPU to allocate to the container. Possible values include `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, and `2.0`. When there's a workload profile specified, there's no such constraint.

~> **Note:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.0` / `2.0` or `0.5` / `1.0`

* `env` - (Optional) One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

~> **Note:** `ephemeral_storage` is currently in preview and not configurable at this time.

* `image` - (Required) The image to use to create the container.

* `liveness_probe` - (Optional) A `liveness_probe` block as detailed below.

* `memory` - (Required) The amount of memory to allocate to the container. Possible values are `0.5Gi`, `1Gi`, `1.5Gi`, `2Gi`, `2.5Gi`, `3Gi`, `3.5Gi` and `4Gi`. When there's a workload profile specified, there's no such constraint.

~> **Note:** `cpu` and `memory` must be specified in `0.25'/'0.5Gi` combination increments. e.g. `1.25` / `2.5Gi` or `0.75` / `1.5Gi`

* `name` - (Required) The name of the container

* `readiness_probe` - (Optional) A `readiness_probe` block as detailed below.

* `startup_probe` - (Optional) A `startup_probe` block as detailed below.

* `volume_mounts` - (Optional) A `volume_mounts` block as detailed below.

---

A `liveness_probe` block supports the following:

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `1` seconds.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are in the range `1` - `240`. Defaults to `10`.

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

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

~> **Note:** This value is ignored if `secret_name` is used

---

A `readiness_probe` block supports the following:

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.

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

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `30`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The value for the host header which should be sent with this probe. If unspecified, the IP Address of the Pod is used as the host header. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

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

* `sub_path` - (Optional) The sub path of the volume to be mounted in the container.
---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ingress` block supports the following:

* `allow_insecure_connections` - (Optional) Should this ingress allow insecure connections?

* `fqdn` - The FQDN of the ingress.

* `external_enabled` - (Optional) Are connections to this Ingress from outside the Container App Environment enabled? Defaults to `false`.

* `ip_security_restriction` - (Optional) One or more `ip_security_restriction` blocks for IP-filtering rules as defined below.

* `target_port` - (Required) The target port on the container for the Ingress traffic.

* `exposed_port` - (Optional) The exposed port on the container for the Ingress traffic.

~> **Note:** `exposed_port` can only be specified when `transport` is set to `tcp`.

* `traffic_weight` - (Required) One or more `traffic_weight` blocks as detailed below.

* `transport` - (Optional) The transport method for the Ingress. Possible values are `auto`, `http`, `http2` and `tcp`. Defaults to `auto`.

~> **Note:** if `transport` is set to `tcp`, `exposed_port` and `target_port` should be set at the same time.

* `client_certificate_mode` - (Optional) The client certificate mode for the Ingress. Possible values are `require`, `accept`, and `ignore`.

---

A `ip_security_restriction` block supports the following:

* `action` - (Required) The IP-filter action. `Allow` or `Deny`.

~> **Note:** The `action` types in an all `ip_security_restriction` blocks must be the same for the `ingress`, mixing `Allow` and `Deny` rules is not currently supported by the service.

* `description` - (Optional) Describe the IP restriction rule that is being sent to the container-app.

* `ip_address_range` - (Required) The incoming IP address or range of IP addresses (in CIDR notation).

* `name` - (Required) Name for the IP restriction rule.

---

A `traffic_weight` block supports the following:

~> **Note:** This block only applies when `revision_mode` is set to `Multiple`.

* `label` - (Optional) The label to apply to the revision as a name prefix for routing traffic.

* `latest_revision` - (Optional) This traffic Weight applies to the latest stable Container Revision. At most only one `traffic_weight` block can have the `latest_revision` set to `true`.

* `revision_suffix` - (Optional) The suffix string to which this `traffic_weight` applies.

~> **Note:** If `latest_revision` is `false`, the `revision_suffix` shall be specified.

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

~> **Note:** The Resource ID must be of a User Assigned Managed identity defined in an `identity` block.

* `password_secret_name` - (Optional) The name of the Secret Reference containing the password value for this user on the Container Registry, `username` must also be supplied.

* `username` - (Optional) The username to use for this Container Registry, `password_secret_name` must also be supplied..

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App.

* `custom_domain_verification_id` - The ID of the Custom Domain Verification for this Container App.

* `ingress` - An `ingress` block as detailed below.

* `latest_revision_fqdn` - The FQDN of the Latest Revision of the Container App.

* `latest_revision_name` - The name of the latest Container Revision.

* `location` - The location this Container App is deployed in. This is the same as the Environment in which it is deployed.

* `outbound_ip_addresses` - A list of the Public IP Addresses which the Container App uses for outbound network access.

---

An `ingress` block exports the following:

* `custom_domain` - One or more `custom_domain` block as detailed below.

---

A `custom_domain` block exports the following:

* `certificate_binding_type` - The Binding type.

* `certificate_id` - The ID of the Container App Environment Certificate.

* `name` - The hostname of the Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App.
* `update` - (Defaults to 30 minutes) Used when updating the Container App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App.

## Import

A Container App can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_app.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.App/containerApps/myContainerApp"
```
