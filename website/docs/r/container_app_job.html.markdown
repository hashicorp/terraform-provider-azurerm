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
    containers {
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Container App Job resource. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Container App Job. Changing this forces a new resource to be created.

* `container_app_environment_id` - (Required) The ID of the Container App Environment in which to create the Container App Job. Changing this forces a new resource to be created.

* `workload_profile_name` - (Optional) The name of the workload profile to use for the Container App Job. Changing this forces a new resource to be created.

* `replica_timeout_in_seconds` - (Required) The maximum number of seconds a replica is allowed to run.

* `replica_retry_limit` - (Optional) The maximum number of times a replica is allowed to retry.

* `secrets` - (Optional) A `secrets` block as defined below.

* `registries` - (Optional) A `registries` block as defined below.

* `manual_trigger_config` - (Optional) A `manual_trigger_config` block as defined below.

* `event_trigger_config` - (Optional) A `event_trigger_config` block as defined below.

* `schedule_trigger_config` - (Optional) A `schedule_trigger_config` block as defined below.

~> ** NOTE **: Only one of `manual_trigger_config`, `event_trigger_config` or `schedule_trigger_config` can be specified.

*  `template` - (Optional) A `template` block as defined below.

* `identity` - (Optional) A `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `template` block supports the following:

* `containers` - (Optional) A `containers` block as defined below.

* `volumes` - (Optional) A `volumes` block as defined below.

---

A `containers` block supports the following:

* `name` - (Required) The name of the container.

* `image` - (Required) The container image name and tag.

* `cpu` - (Required) The CPU requirement of the container.

* `memory` - (Required) The memory requirement of the container.

* `command` - (Optional) The command to execute within the container.

* `args` - (Optional) The arguments to the command.

* `env` - (Optional) A `env` block as defined below.

* `liveness_probes` - (Optional) A `live_probes` block as defined below.

* `readiness_probes` - (Optional) A `readiness_probes` block as defined below.

* `startup_probes` - (Optional) A `startup_probes` block as defined below.

* `volume_mounts` - (Optional) A `volume_mounts` block as defined below.

---

A `env` block supports the following:

* `name` - (Required) The name of the environment variable.

* `value` - (Optional) The value of the environment variable.

* `secret_name` - (Optional) Name of the Container App secret from which to pull the environment variable value.

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

* `name` - (Required) The name of the volume to mount. This must match the name of a volume defined in the `volumes` block.

* `path` - (Required) The path within the container at which the volume should be mounted. Must not contain `:`.

---

A `volumes` block supports the following:

* `name` - (Optional) The name of the volume.

* `storage_type` - (Optional) The type of storage to use for the volume. Possible values are `AzureFile`, `EmptyDir` and `Secret`.

* `storage_name` - (Optional) The name of the storage to use for the volume.

---

A `secrets` block supports the following:

* `name` - (required) Name of the secret.

* `value` - (required) Value of the secret.

---

A `registries` block supports the following:

* `identity` - (Optional) A Managed Identity to use to authenticate with Azure Container Registry.

* `username` - (Optional) The username to use to authenticate with Azure Container Registry.

* `password_secret_name` - (Optional) The name of the Secret that contains the registry login password.

* `server` - (Optional) The URL of the Azure Container Registry server.

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

* `polling_interval` - (Optional) Interval to check each event source in seconds.

* `rules` - (Optional) A `rules` block as defined below.

---

A `rules` block supports the following:

* `name` - (Optional) Name of the scale rule.

* `type` - (Optional) Type of the scale rule.

* `metadata` - (Optional) Metadata properties to describe the scale rule.

* `auth` - (Optional) A `auth` block as defined below.

---

A `auth` block supports the following:

* `secret_ref` - (Optional) Name of the secret from which to pull the auth params.

* `trigger_parameter` - (Optional) Trigger Parameter that uses the secret.

---

A `identity` block supports the following:

* `type` - (Optional) The type of identity used for the Container App Job. Possible values are `SystemAssigned` and `None`. Defaults to `None`.

* `identity_ids` - (Optional) A list of Managed Identity IDs to assign to the Container App Job.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Job.

* `outbound_ip_addresses` - A list of the Public IP Addresses which the Container App uses for outbound network access.

* `event_stream_endpoint` - The endpoint for the Container App Job event stream.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to `30 minutes`) Used when creating the Container App Job.
* `update` - (Defaults to `30 minutes`) Used when updating the Container App Job.
* `read` - (Defaults to `5 minutes`) Used when retrieving the Container App Job.
* `delete` - (Defaults to `30 minutes`) Used when deleting the Container App Job.

## Import

A Container App Job can be imported using the resource id, e.g.

```shell
terraform import azurerm_container_app_job.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.App/jobs/example-container-app-job"
```
