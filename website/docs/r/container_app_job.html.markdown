---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_job"
description: |-
  Manages a Container App Job.
---

# azurerm_container_app_job

Manages a Container App Job.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-log-analytics-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "example" {
  name                       = "example-container-app-environment"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_container_app_job" "example" {
  name                         = "example-container-app-job"
  location                     = azurerm_resource_group.example.location
  resource_group_name          = azurerm_resource_group.example.name
  container_app_environment_id = azurerm_container_app_environment.example.id

  replica_timeout_in_seconds = 10
  replica_retry_limit        = 10
  manual_trigger_config {
    parallelism              = 4
    replica_completion_count = 1
  }

  template {
    container {
      image = "repo/testcontainerAppsJob0:v1"
      name  = "testcontainerappsjob0"
      readiness_probe {
        transport = "HTTP"
        port      = 5000
      }

      liveness_probe {
        transport = "HTTP"
        port      = 5000
        path      = "/health"

        header {
          name  = "Cache-Control"
          value = "no-cache"
        }

        initial_delay           = 5
        interval_seconds        = 20
        timeout                 = 2
        failure_count_threshold = 1
      }
      startup_probe {
        transport = "TCP"
        port      = 5000
      }

      cpu    = 0.5
      memory = "1Gi"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container App Job resource. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container App Job. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The ID of the Container App Environment in which to create the Container App Job. Changing this forces a new resource to be created.

* `template` - (Required) A `template` block as defined below.

* `replica_timeout_in_seconds` - (Required) The maximum number of seconds a replica is allowed to run.

* `workload_profile_name` - (Optional) The name of the workload profile to use for the Container App Job.

* `replica_retry_limit` - (Optional) The maximum number of times a replica is allowed to retry.

* `secret` - (Optional) One or more `secret` blocks as defined below.

* `registry` - (Optional) One or more `registry` blocks as defined below.

* `manual_trigger_config` - (Optional) A `manual_trigger_config` block as defined below. Changing this forces a new resource to be created.

* `event_trigger_config` - (Optional) A `event_trigger_config` block as defined below. Changing this forces a new resource to be created.

* `schedule_trigger_config` - (Optional) A `schedule_trigger_config` block as defined below. Changing this forces a new resource to be created.

~> **Note:** Only one of `manual_trigger_config`, `event_trigger_config` or `schedule_trigger_config` can be specified.

* `identity` - (Optional) A `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `template` block supports the following:

* `container` - (Required) A `container` block as defined below.

* `init_container` - (Optional) A `init_container` block as defined below.

* `volume` - (Optional) A `volume` block as defined below.

---

A `container` block supports the following:

* `name` - (Required) The name of the container.

* `cpu` - (Required) The amount of vCPU to allocate to the container.

~> **Note:** When using a Consumption plan, the `cpu` and `memory` properties must add up to one of the combinations found in the Microsoft provided documentation, for more information see [vCPU and memory allocation requirements](https://learn.microsoft.com/azure/container-apps/containers#allocations)

* `memory` - (Required) The amount of memory to allocate to the container.

~> **Note:** When using a Consumption plan, the `cpu` and `memory` properties must add up to one of the combinations found in the Microsoft provided documentation, for more information see [vCPU and memory allocation requirements](https://learn.microsoft.com/azure/container-apps/containers#allocations)

* `image` - (Required) The image to use to create the container.

* `args` - (Optional) A list of extra arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `env` - (Optional) One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

~> **Note:** `ephemeral_storage` is currently in preview and not configurable at this time.

* `liveness_probe` - (Optional) A `liveness_probe` block as detailed below.

* `readiness_probe` - (Optional) A `readiness_probe` block as detailed below.

* `startup_probe` - (Optional) A `startup_probe` block as detailed below.

* `volume_mounts` - (Optional) A `volume_mounts` block as detailed below.

---

An `init_container` block supports:

* `name` - (Required) The name of the container.

* `cpu` - (Optional) The amount of vCPU to allocate to the container.

~> **Note:** When using a Consumption plan, the `cpu` and `memory` properties must add up to one of the combinations found in the Microsoft provided documentation, for more information see [vCPU and memory allocation requirements](https://learn.microsoft.com/azure/container-apps/containers#allocations)

* `memory` - (Optional) The amount of memory to allocate to the container.

~> **Note:** When using a Consumption plan, the `cpu` and `memory` properties must add up to one of the combinations found in the Microsoft provided documentation, for more information see [vCPU and memory allocation requirements](https://learn.microsoft.com/azure/container-apps/containers#allocations)

* `image` - (Required) The image to use to create the container.

* `args` - (Optional) A list of extra arguments to pass to the container.

* `command` - (Optional) A command to pass to the container to override the default. This is provided as a list of command line elements without spaces.

* `env` - (Optional) One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the Container App.

~> **Note:** `ephemeral_storage` is currently in preview and not configurable at this time.

* `volume_mounts` - (Optional) A `volume_mounts` block as detailed below.

---

A `env` block supports the following:

* `name` - (Required) The name of the environment variable.

* `value` - (Optional) The value of the environment variable.

* `secret_name` - (Optional) Name of the Container App secret from which to pull the environment variable value.

---

A `liveness_probe` block supports the following:

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The time in seconds to wait after the container has started before the probe is started.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are in the range `1` - `240`. Defaults to `10`.

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

A `readiness_probe` block supports the following:

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The probe hostname. Defaults to the pod IP address. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - (Optional) The URI to use for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `success_count_threshold` - (Optional) The number of consecutive successful responses required to consider this probe as successful. Possible values are between `1` and `10`. Defaults to `3`.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

A `startup_probe` block supports the following:

* `port` - (Required) The port number on which to connect. Possible values are between `1` and `65535`.

* `transport` - (Required) Type of probe. Possible values are `TCP`, `HTTP`, and `HTTPS`.

* `failure_count_threshold` - (Optional) The number of consecutive failures required to consider this probe as failed. Possible values are between `1` and `10`. Defaults to `3`.

* `header` - (Optional) A `header` block as detailed below.

* `host` - (Optional) The value for the host header which should be sent with this probe. If unspecified, the IP Address of the Pod is used as the host header. Setting a value for `Host` in `headers` can be used to override this for `HTTP` and `HTTPS` type probes.

* `initial_delay` - (Optional) The number of seconds elapsed after the container has started before the probe is initiated. Possible values are between `0` and `60`. Defaults to `0` seconds.

* `interval_seconds` - (Optional) How often, in seconds, the probe should run. Possible values are between `1` and `240`. Defaults to `10`

* `path` - (Optional) The URI to use with the `host` for http type probes. Not valid for `TCP` type probes. Defaults to `/`.

* `timeout` - (Optional) Time in seconds after which the probe times out. Possible values are in the range `1` - `240`. Defaults to `1`.

---

A `header` block supports the following:

* `name` - (Required) The HTTP Header Name.

* `value` - (Required) The HTTP Header value.

---

A `volume_mounts` block supports the following:

* `name` - (Required) The name of the volume to mount. This must match the name of a volume defined in the `volume` block.

* `path` - (Required) The path within the container at which the volume should be mounted. Must not contain `:`.

* `sub_path` - (Optional) The sub path of the volume to be mounted in the container.

---

A `volume` block supports the following:

* `name` - (Required) The name of the volume.

* `storage_type` - (Optional) The type of storage to use for the volume. Possible values are `AzureFile`, `EmptyDir`, `NfsAzureFile` and `Secret`. Defaults to `EmptyDir`.

* `storage_name` - (Optional) The name of the storage to use for the volume.

* `mount_options` - (Optional) Mount options used while mounting the AzureFile. Must be a comma-separated string e.g. `dir_mode=0751,file_mode=0751`.

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

A `registry` block supports the following:

* `identity` - (Optional) A Managed Identity to use to authenticate with Azure Container Registry.

* `username` - (Optional) The username to use to authenticate with Azure Container Registry.

* `password_secret_name` - (Optional) The name of the Secret that contains the registry login password.

* `server` - (Required) The URL of the Azure Container Registry server.

---

A `manual_trigger_config` block supports the following:

* `parallelism` - (Optional) Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - (Optional) Minimum number of successful replica completions before overall job completion.

---

A `event_trigger_config` block supports the following:

* `parallelism` - (Optional) Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - (Optional) Minimum number of successful replica completions before overall job completion.

* `scale` - (Optional) A `scale` block as defined below.

---

A `schedule_trigger_config` block supports the following:

* `cron_expression` - (Required) Cron formatted repeating schedule of a Cron Job.

* `parallelism` - (Optional) Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - (Optional) Minimum number of successful replica completions before overall job completion.

---

A `scale` block supports the following:

* `max_executions` - (Optional) Maximum number of job executions that are created for a trigger.

* `min_executions` - (Optional) Minimum number of job executions that are created for a trigger.

* `polling_interval_in_seconds` - (Optional) Interval to check each event source in seconds.

* `rules` - (Optional) A `rules` block as defined below.

---

A `rules` block supports the following:

* `name` - (Required) Name of the scale rule.

* `custom_rule_type` - (Required) Type of the scale rule. Possible values are `activemq`, `artemis-queue`, `kafka`, `pulsar`, `aws-cloudwatch`, `aws-dynamodb`, `aws-dynamodb-streams`, `aws-kinesis-stream`, `aws-sqs-queue`, `azure-app-insights`, `azure-blob`, `azure-data-explorer`, `azure-eventhub`, `azure-log-analytics`, `azure-monitor`, `azure-pipelines`, `azure-servicebus`, `azure-queue`, `cassandra`, `cpu`, `cron`, `datadog`, `elasticsearch`, `external`, `external-push`, `gcp-stackdriver`, `gcp-storage`, `gcp-pubsub`, `graphite`, `http`, `huawei-cloudeye`, `ibmmq`, `influxdb`, `kubernetes-workload`, `liiklus`, `memory`, `metrics-api`, `mongodb`, `mssql`, `mysql`, `nats-jetstream`, `stan`, `tcp`, `new-relic`, `openstack-metric`, `openstack-swift`, `postgresql`, `predictkube`, `prometheus`, `rabbitmq`, `redis`, `redis-cluster`, `redis-sentinel`, `redis-streams`, `redis-cluster-streams`, `redis-sentinel-streams`, `selenium-grid`, `solace-event-queue` and `github-runner`.

* `metadata` - (Required) Metadata properties to describe the scale rule.

* `authentication` - (Optional) A `authentication` block as defined below.

---

A `authentication` block supports the following:

* `secret_name` - (Required) Name of the secret from which to pull the auth params.

* `trigger_parameter` - (Required) Trigger Parameter that uses the secret.

---

A `identity` block supports the following:

* `type` - (Required) The type of identity used for the Container App Job. Possible values are `SystemAssigned`, `UserAssigned` and `None`. Defaults to `None`.

* `identity_ids` - (Optional) A list of Managed Identity IDs to assign to the Container App Job.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Job.

* `outbound_ip_addresses` - A list of the Public IP Addresses which the Container App uses for outbound network access.

* `event_stream_endpoint` - The endpoint for the Container App Job event stream.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container App Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Job.
* `update` - (Defaults to 30 minutes) Used when updating the Container App Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container App Job.

## Import

A Container App Job can be imported using the resource id, e.g.

```shell
terraform import azurerm_container_app_job.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.App/jobs/example-container-app-job"
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
