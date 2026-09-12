---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_job"
description: |-
  Gets information about an existing Container App Job.
---

# Data Source: azurerm_container_app_job

Use this data source to access information about an existing Container App Job.

## Example Usage

```hcl
data "azurerm_container_app_job" "example" {
  name                = "example-container-app-job"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container App Job.

* `resource_group_name` - (Required) The name of the Resource Group where this Container App Job exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Job.

* `location` - The Azure Region where the Container App Job exists.

* `container_app_environment_id` - The ID of the Container App Environment this Container App Job is linked to.

* `event_stream_endpoint` - The endpoint used to stream events from the Container App Job.

* `event_trigger_config` - An `event_trigger_config` block as detailed below.

* `identity` - An `identity` block as detailed below.

* `manual_trigger_config` - A `manual_trigger_config` block as detailed below.

* `outbound_ip_addresses` - A list of the Public IP Addresses which the Container App Job uses for outbound network access.

* `registry` - One or more `registry` blocks as detailed below.

* `replica_retry_limit` - The maximum number of times a replica is allowed to retry.

* `replica_timeout_in_seconds` - The maximum time in seconds that a replica is allowed to run.

* `schedule_trigger_config` - A `schedule_trigger_config` block as detailed below.

* `secret` - One or more `secret` blocks as detailed below.

* `template` - A `template` block as detailed below.

* `workload_profile_name` - The name of the Workload Profile in the Container App Environment in which this Container App Job is running.

* `tags` - A mapping of tags assigned to the Container App Job.

---

A `template` block exports the following:

* `container` - One or more `container` blocks as detailed below.

* `init_container` - One or more `init_container` blocks as detailed below.

* `volume` - One or more `volume` blocks as detailed below.

---

A `container` block exports the following:

* `args` - A list of extra arguments passed to the container.

* `command` - A command passed to the container to override the default.

* `cpu` - The amount of vCPU allocated to the container.

* `env` - One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the container.

* `image` - The image used to create the container.

* `liveness_probe` - A `liveness_probe` block as detailed below.

* `memory` - The amount of memory allocated to the container.

* `name` - The name of the container.

* `readiness_probe` - A `readiness_probe` block as detailed below.

* `startup_probe` - A `startup_probe` block as detailed below.

* `volume_mounts` - A `volume_mounts` block as detailed below.

---

An `init_container` block exports the following:

* `args` - A list of extra arguments passed to the container.

* `command` - A command passed to the container to override the default.

* `cpu` - The amount of vCPU allocated to the container.

* `env` - One or more `env` blocks as detailed below.

* `ephemeral_storage` - The amount of ephemeral storage available to the container.

* `image` - The image used to create the container.

* `memory` - The amount of memory allocated to the container.

* `name` - The name of the container.

* `volume_mounts` - A `volume_mounts` block as detailed below.

---

An `env` block exports the following:

* `name` - The name of the environment variable.

* `secret_name` - The name of the secret that contains the value for this environment variable.

* `value` - The value of the environment variable.

---

A `liveness_probe`, `readiness_probe` or `startup_probe` block exports the following:

* `failure_count_threshold` - The number of consecutive failures required to consider this probe as failed.

* `header` - A `header` block as detailed below.

* `host` - The probe hostname.

* `initial_delay` - The number of seconds elapsed after the container has started before the probe is initiated.

* `interval_seconds` - How often, in seconds, the probe should run.

* `path` - The URI used for HTTP type probes.

* `port` - The port number on which to connect.

* `success_count_threshold` - The number of consecutive successful responses required to consider this probe as successful. Only valid for `readiness_probe`.

* `timeout` - Time in seconds after which the probe times out.

* `transport` - The type of probe.

* `termination_grace_period_seconds` - The time in seconds after the container is sent the termination signal before the process is forcibly killed.

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

A `volume` block exports the following:

* `name` - The name of the volume.

* `storage_name` - The name of the storage to use for the volume.

* `storage_type` - The type of storage volume.

* `mount_options` - Mount options used while mounting the volume.

---

An `event_trigger_config` block exports the following:

* `parallelism` - Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - Minimum number of successful replica completions before overall job completion.

* `scale` - A `scale` block as detailed below.

---

A `scale` block exports the following:

* `max_executions` - Maximum number of job executions that are created for a trigger.

* `min_executions` - Minimum number of job executions that are created for a trigger.

* `polling_interval_in_seconds` - Interval to check each event source in seconds.

* `rules` - One or more `rules` blocks as detailed below.

---

A `rules` block exports the following:

* `name` - Name of the scale rule.

* `custom_rule_type` - Type of the custom scale rule.

* `metadata` - Metadata properties to describe the custom scale rule.

* `authentication` - One or more `authentication` blocks as detailed below.

---

An `authentication` block exports the following:

* `secret_name` - The name of the Container App secret to use for this Scale Rule Authentication.

* `trigger_parameter` - The Trigger Parameter name to use the supply the value retrieved from the `secret_name`.

---

A `manual_trigger_config` block exports the following:

* `parallelism` - Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - Minimum number of successful replica completions before overall job completion.

---

A `schedule_trigger_config` block exports the following:

* `cron_expression` - Cron formatted repeating schedule of a Cron Job.

* `parallelism` - Number of parallel replicas of a job that can run at a given time.

* `replica_completion_count` - Minimum number of successful replica completions before overall job completion.

---

A `registry` block exports the following:

* `server` - The hostname for the Container Registry.

* `identity` - The Resource ID for the User Assigned Managed identity used to authenticate with the Container Registry.

* `password_secret_name` - The name of the Secret reference that contains the password value for the Container Registry user.

* `username` - The username used to authenticate with the Container Registry.

---

A `secret` block exports the following:

* `identity` - The identity used for accessing the Key Vault.

* `key_vault_secret_id` - The ID of a Key Vault secret.

* `name` - The secret name.

* `value` - The value for this secret.

---

An `identity` block exports the following:

* `type` - The type of Managed Identity assigned to this Container App Job.

* `identity_ids` - A list of the User Assigned Managed Identity IDs assigned to this Container App Job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Job.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
