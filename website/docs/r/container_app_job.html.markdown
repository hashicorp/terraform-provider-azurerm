---
subcategory: "Container App"
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

  configuration {
    trigger_type        = "Manual"
    replica_timeout     = 10
    replica_retry_limit = 10
    manual_trigger_config {
      parallelism              = 4
      replica_completion_count = 1
    }
  }

  template {
    containers {
      image = "repo/testcontainerAppsJob0:v1"
      name  = "testcontainerappsjob0"
      probes {
        http_get {
          http_headers {
            name  = "testheader"
            value = "testvalue"
          }
          path = "/testpath"
          port = 8080
        }
        initial_delay_seconds = 10
        period_seconds        = 10
        type                  = "Liveness"
      }
      resources {
        cpu    = 0.5
        memory = "1Gi"
      }
    }

    init_containers {
      args    = ["testarg"]
      command = ["testcommand"]
      image   = "repo/testcontainerAppsJob0:v1"
      name    = "testcontainerappsjob0"
      resources {
        cpu    = 0.5
        memory = "1Gi"
      }
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

*  `template` - (Optional) A `template` block as defined below.

* `configuration` - (Optional) A `configuration` block as defined below.

* `identity` - (Optional) A `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `template` block supports the following:

* `containers` - (Optional) A `containers` block as defined below.

* `init_containers` - (Optional) A `init_containers` block as defined below.

* `volumes` - (Optional) A `volumes` block as defined below.

---

A `containers` block supports the following:

* `name` - (Required) The name of the container.

* `image` - (Required) The container image name and tag.

* `command` - (Optional) The command to execute within the container.

* `args` - (Optional) The arguments to the command.

* `env` - (Optional) A `env` block as defined below.

* `probes` - (Optional) A `probes` block as defined below.

* `resources` - (Optional) A `resources` block as defined below.

* `volume_mounts` - (Optional) A `volume_mounts` block as defined below.

---

A `init_containers` block supports the following:

* `name` - (Required) The name of the container.

* `image` - (Required) The container image name and tag.

* `command` - (Optional) The command to execute within the container.

* `args` - (Optional) The arguments to the command.

* `env` - (Optional) A `env` block as defined below.

* `probes` - (Optional) A `probes` block as defined below.

* `resources` - (Optional) A `resources` block as defined below.

* `volume_mounts` - (Optional) A `volume_mounts` block as defined below.

---

A `env` block supports the following:

* `name` - (Optional) The name of the environment variable.

* `value` - (Optional) The value of the environment variable.

* `secret_ref` - (Optional) Name of the Container App secret from which to pull the environment variable value.

---

A `probes` block supports the following:

* `type` - (Required) The type of probe.

* `http_get` - (Optional) A `http_get` block as defined below.

* `failure_threshold` - (Optional) The number of times the probe can fail before the container is restarted. Defaults to 3.

* `initial_delay_seconds` - (Optional) The number of seconds after the container has started before the probe is initiated.

* `period_seconds` - (Optional) The number of seconds between probe checks. Defaults to 10.

* `success_threshold` - (Optional) The number of times the probe must succeed before the container is marked as running. Defaults to 1.

* `timeout_seconds` - (Optional) The number of seconds after which the probe times out.

* `termination_grace_period_seconds` - (Optional) The number of seconds after which the container is terminated if the probe failed.

* `tcp_socket` - (Optional) A `tcp_socket` block as defined below.

---

A `http_get` block supports the following:

* `port` - (Required) Name or number of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.

* `path` - (Optional) Path to access on the HTTP server.

* `scheme` - (Optional) The scheme to use for probing. Possible values are `HTTP` and `HTTPS`. 

* `http_headers` - (Optional) A `http_headers` block as defined below.

---

A `http_headers` block supports the following:

* `name` - (Required) The name of the header.

* `value` - (Required) The value of the header.

---

A `tcp_socket` block supports the following:

* `port` - (Required) TCPSocket specifies an action involving a TCP port.

* `host` - (Optional) Host name to connect to, defaults to the pod IP.

---

A `resources` block supports the following:

* `cpu` - (Required) Required CPU in cores. e.g. `0.5`.

* `memory` - (Required) Required memory in GB. e.g. `1Gi`.

---

A `volume_mounts` block supports the following:

* `volume_name` - (Optional) The name of the volume to mount. This must match the name of a volume defined in the `volumes` block.

* `mount_path` - (Optional) The path within the container at which the volume should be mounted. Must not contain `:`.

* `sub_path` - (Optional) Path within the volume from which the container's volume should be mounted.

---

A `volumes` block supports the following:

* `mount_options` - (Optional) The mount options for the volume, e.g. `["ro", "soft"]`.

* `name` - (Optional) The name of the volume.

* `secrets` - (Optional) A `secrets` block as defined below.

* `storage_type` - (Optional) The type of storage to use for the volume. Possible values are `AzureFile`, `EmptyDir` and `Secret`.

* `storage_name` - (Optional) The name of the storage to use for the volume.

---

A `secrets` block supports the following:

* `path` - (Optional) Path to project secret to. If no path is provided, path defaults to name of secret listed in secretRef.

* `secret_ref` - (Optional) Name of the Container App secret from which to pull the secret value.

---

A `configuration` block supports the following:

* `trigger_type` - (Required) The type of trigger for the Container App Job. Possible values are `Manual`, `Event` and `Schedule`.

* `replica_timeout` - (Required) The maximum number of seconds a replica is allowed to run.

* `replica_retry_limit` - (Optional) The maximum number of times a replica is allowed to retry.

* `secret` - (Optional) A `secret` block as defined below.

* `registries` - (Optional) A `registries` block as defined below.

* `manual_trigger_config` - (Optional) A `manual_trigger_config` block as defined below.

* `event_trigger_config` - (Optional) A `event_trigger_config` block as defined below.

* `schedule_trigger_config` - (Optional) A `schedule_trigger_config` block as defined below.

~> ** NOTE **: Only one of `manual_trigger_config`, `event_trigger_config` or `schedule_trigger_config` can be specified.

---

A `secret` block supports the following:

* `identity` - (Optional) Resource ID of a managed identity to authenticate with Azure Key Vault, or System to use a system-assigned identity.

* `key_vault_uri` - (Optional) Azure Key Vault URL pointing to the secret referenced by the container app.

* `name` - (Optional) Name of the secret.

* `value` - (Optional) Value of the secret.

---

A `registries` block supports the following:

* `identity` - (Optional) A Managed Identity to use to authenticate with Azure Container Registry.

* `username` - (Optional) The username to use to authenticate with Azure Container Registry.

* `password_secret_ref` - (Optional) The name of the Secret that contains the registry login password.

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

* `trigger_parameter` - (Optional) Trigger Parameter that uses the secret.\

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
