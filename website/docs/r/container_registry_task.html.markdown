---
subcategory: "Container"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_registry_task"
description: |-
  Manages a Container Registry Task.
---

# azurerm_container_registry_task

Manages a Container Registry Task.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}
resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Basic"
}
resource "azurerm_container_registry_task" "example" {
  name                  = "example-task"
  container_registry_id = azurerm_container_registry.example.id
  platform {
    os = "Linux"
  }
  docker_step {
    dockerfile_path      = "Dockerfile"
    context_path         = "https://github.com/<username>/<repository>#<branch>:<folder>"
    context_access_token = "<github personal access token>"
    image_names          = ["helloworld:{{.Run.ID}}"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Container Registry Task. Changing this forces a new Container Registry Task to be created.

* `container_registry_id` - (Required) The ID of the Container Registry that this Container Registry Task resides in. Changing this forces a new Container Registry Task to be created.

---

* `agent_pool_name` - (Optional) The name of the dedicated Container Registry Agent Pool for this Container Registry Task.

* `agent_setting` - (Optional) A `agent_setting` block as defined below.

~> **Note:** Only one of `agent_pool_name` and `agent_setting` can be specified.

* `enabled` - (Optional) Should this Container Registry Task be enabled? Defaults to `true`.

* `identity` - (Optional) An `identity` block as defined below.

* `platform` - (Optional) A `platform` block as defined below.

~> **Note:** The `platform` is required for non-system task (when `is_system_task` is set to `false`).

* `docker_step` - (Optional) A `docker_step` block as defined below.

* `encoded_step` - (Optional) A `encoded_step` block as defined below.

* `file_step` - (Optional) A `file_step` block as defined below.

~> **Note:** For non-system task (when `is_system_task` is set to `false`), one and only one of the `docker_step`, `encoded_step` and `file_step` should be specified.

* `base_image_trigger` - (Optional) A `base_image_trigger` block as defined below.

* `source_trigger` - (Optional) One or more `source_trigger` blocks as defined below.

* `timer_trigger` - (Optional) One or more `timer_trigger` blocks as defined below.

* `is_system_task` - (Optional) Whether this Container Registry Task is a system task. Changing this forces a new Container Registry Task to be created. Defaults to `false`.

~> **Note:** For system task, the `name` has to be set as `quicktask`. And the following properties can't be specified: `docker_step`, `encoded_step`, `file_step`, `platform`, `base_image_trigger`, `source_trigger`, `timer_trigger`.

* `log_template` - (Optional) The template that describes the run log artifact.

* `registry_credential` - (Optional) One `registry_credential` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Container Registry Task.

* `timeout_in_seconds` - (Optional) The timeout of this Container Registry Task in seconds. The valid range lies from 300 to 28800. Defaults to `3600`.

---

A `agent_setting` block supports the following:

* `cpu` - (Required) The number of cores required for the Container Registry Task. Possible value is `2`.

---

A `authentication` block supports the following:

* `token` - (Required) The access token used to access the source control provider.

* `token_type` - (Required) The type of the token. Possible values are `PAT` (personal access token) and `OAuth`.

* `expire_in_seconds` - (Optional) Time in seconds that the token remains valid.

* `refresh_token` - (Optional) The refresh token used to refresh the access token.

* `scope` - (Optional) The scope of the access token.

---

A `base_image_trigger` block supports the following:

* `name` - (Required) The name which should be used for this trigger.

* `type` - (Required) The type of the trigger. Possible values are `All` and `Runtime`.

* `enabled` - (Optional) Should the trigger be enabled? Defaults to `true`.

* `update_trigger_endpoint` - (Optional) The endpoint URL for receiving the trigger.

* `update_trigger_payload_type` - (Optional) Type of payload body for the trigger. Possible values are `Default` and `Token`.

---

A `custom` block supports the following:

* `login_server` - (Required) The login server of the custom Container Registry.

* `identity` - (Optional) The managed identity assigned to this custom credential. For user assigned identity, the value is the client ID of the identity. For system assigned identity, the value is `[system]`.

* `password` - (Optional) The password for logging into the custom Container Registry. It can be either a plain text of password, or a Keyvault Secret ID.

* `username` - (Optional) The username for logging into the custom Container Registry. It can be either a plain text of username, or a Keyvault Secret ID.

---

A `docker_step` block supports the following:

* `context_access_token` - (Required) The token (Git PAT or SAS token of storage account blob) associated with the context for this step.

* `context_path` - (Required) The URL (absolute or relative) of the source context for this step. If the context is an url you can reference a specific branch or folder via `#branch:folder`.

* `dockerfile_path` - (Required) The Dockerfile path relative to the source context.

* `arguments` - (Optional) Specifies a map of arguments to be used when executing this step.

* `image_names` - (Optional) Specifies a list of fully qualified image names including the repository and tag.

* `cache_enabled` - (Optional) Should the image cache be enabled? Defaults to `true`.

* `push_enabled` - (Optional) Should the image built be pushed to the registry or not? Defaults to `true`.

* `secret_arguments` - (Optional) Specifies a map of *secret* arguments to be used when executing this step.

* `target` - (Optional) The name of the target build stage for the docker build.

---

A `encoded_step` block supports the following:

* `task_content` - (Required) The (optionally base64 encoded) content of the build template.

* `context_access_token` - (Optional) The token (Git PAT or SAS token of storage account blob) associated with the context for this step.

* `context_path` - (Optional) The URL (absolute or relative) of the source context for this step.

* `secret_values` - (Optional) Specifies a map of secret values that can be passed when running a task.

* `value_content` - (Optional) The (optionally base64 encoded) content of the build parameters.

* `values` - (Optional) Specifies a map of values that can be passed when running a task.

---

A `file_step` block supports the following:

* `task_file_path` - (Required) The task template file path relative to the source context.

* `context_access_token` - (Optional) The token (Git PAT or SAS token of storage account blob) associated with the context for this step.

* `context_path` - (Optional) The URL (absolute or relative) of the source context for this step.

* `secret_values` - (Optional) Specifies a map of secret values that can be passed when running a task.

* `value_file_path` - (Optional) The parameters file path relative to the source context.

* `values` - (Optional) Specifies a map of values that can be passed when running a task.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Container Registry Task. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Container Registry Task.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `platform` block supports the following:

* `os` - (Required) The operating system type required for the task. Possible values are `Windows` and `Linux`.

* `architecture` - (Optional) The OS architecture. Possible values are `amd64`, `x86`, `386`, `arm` and `arm64`.

* `variant` - (Optional) The variant of the CPU. Possible values are `v6`, `v7`, `v8`.

---

A `registry_credential` block supports the following:

* `source` - (Optional) One `source` block as defined below.

* `custom` - (Optional) One or more `custom` blocks as defined above.

---

A `source` block supports the following:

* `login_mode` - (Required) The login mode for the source registry. Possible values are `None` and `Default`.

---

A `source_trigger` block supports the following:

* `name` - (Required) The name which should be used for this trigger.

* `events` - (Required) Specifies a list of source events corresponding to the trigger. Possible values are `commit` and `pullrequest`.

* `repository_url` - (Required) The full URL to the source code repository.

* `source_type` - (Required) The type of the source control service. Possible values are `Github` and `VisualStudioTeamService`.

* `authentication` - (Optional) A `authentication` block as defined above.

* `branch` - (Optional) The branch name of the source code.

* `enabled` - (Optional) Should the trigger be enabled? Defaults to `true`.

---

A `timer_trigger` block supports the following:

* `name` - (Required) The name which should be used for this trigger.

* `schedule` - (Required) The CRON expression for the task schedule.

* `enabled` - (Optional) Should the trigger be enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container Registry Task.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Container Registry Task.
* `read` - (Defaults to 5 minutes) Used when retrieving the Container Registry Task.
* `update` - (Defaults to 30 minutes) Used when updating the Container Registry Task.
* `delete` - (Defaults to 30 minutes) Used when deleting the Container Registry Task.

## Import

Container Registry Tasks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_container_registry_task.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.ContainerRegistry/registries/registry1/tasks/task1
```
