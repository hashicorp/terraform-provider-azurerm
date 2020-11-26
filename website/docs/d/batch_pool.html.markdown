---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
description: |-
  Get information about an existing Azure Batch pool.

---

# Data source: azurerm_batch_pool

Use this data source to access information about an existing Batch pool

## Example Usage

```hcl
data "azurerm_batch_pool" "example" {
  name                = "testbatchpool"
  account_name        = "testbatchaccount"
  resource_group_name = "test"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The Batch pool ID.

* `name` - The name of the Batch pool.

* `account_name` - The name of the Batch account.

* `node_agent_sku_id` - The Sku of the node agents in the Batch pool.

* `vm_size` - The size of the VM created in the Batch pool.

* `fixed_scale` - A `fixed_scale` block that describes the scale settings when using fixed scale.

* `auto_scale` - A `auto_scale` block that describes the scale settings when using auto scale.

* `storage_image_reference` - The reference of the storage image used by the nodes in the Batch pool.

* `start_task` - A `start_task` block that describes the start task settings for the Batch pool.

* `max_tasks_per_node` - The maximum number of tasks that can run concurrently on a single compute node in the pool.

* `certificate` - One or more `certificate` blocks that describe the certificates installed on each compute node in the pool.

* `container_configuration` - The container configuration used in the pool's VMs.

---

A `fixed_scale` block exports the following:

* `target_dedicated_nodes` - The number of nodes in the Batch pool.

* `target_low_priority_nodes` - The number of low priority nodes in the Batch pool.

* `resize_timeout` - The timeout for resize operations.

---

A `auto_scale` block exports the following:

* `evaluation_interval` - The interval to wait before evaluating if the pool needs to be scaled.

* `formula` - The autoscale formula that needs to be used for scaling the Batch pool.

---

A `start_task` block exports the following:

* `command_line` - The command line executed by the start task.

* `max_task_retry_count` - The number of retry count.

* `wait_for_success` - A flag that indicates if the Batch pool should wait for the start task to be completed.

* `environment` - A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - A `user_identity` block that describes the user identity under which the start task runs.

* `resource_file` - One or more `resource_file` blocks that describe the files to be downloaded to a compute node.

---

A `user_identity` block exports the following:

* `user_name` - The username to be used by the Batch pool start task.

* `auto_user` - A `auto_user` block that describes the user identity under which the start task runs.

---

A `auto_user` block exports the following:

* `elevation_level` - The elevation level of the user identity under which the start task runs.

* `scope` - The scope of the user identity under which the start task runs.

---

A `certificate` block exports the following:

* `id` - The fully qualified ID of the certificate installed on the pool.

* `store_location` - The location of the certificate store on the compute node into which the certificate is installed, either `CurrentUser` or `LocalMachine`.

-> **NOTE:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of 'remoteUser', a 'certs' directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

* `store_name` - The name of the certificate store on the compute node into which the certificate is installed.

-> **NOTE:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference).

* `visibility` - Which user accounts on the compute node have access to the private data of the certificate.

---

A `resource_file` block exports the following:

* `auto_storage_container_name` - The storage container name in the auto storage account.

* `blob_prefix` - The blob prefix used when downloading blobs from an Azure Storage container.

* `file_mode` - The file permission mode attribute represented as a string in octal format (e.g. `"0644"`).

* `file_path` - The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified.

* `http_url` - The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access.

* `storage_container_url` - The URL of the blob container within Azure Blob Storage.

---

A `container_configuration` block exports the following:

* `type` - The type of container configuration.

* `container_image_names` - A list of container image names to use, as would be specified by `docker pull`.

* `container_registries` - Additional container registries from which container images can be pulled by the pool's VMs.

---

A `container_registries` block exports the following:

* `registry_server` - The container registry URL. The default is "docker.io".

* `user_name` - The user name to log into the registry server.

* `password` - The password to log into the registry server.

---

A `network_configuration` block exports the following:

* `subnet_id` - The ARM resource identifier of the virtual network subnet which the compute nodes of the pool are joined too.

* `endpoint_configuration` - The inbound NAT pools that are used to address specific ports on the individual compute node externally.

---

A `endpoint_configuration` block exports the following:

* `name` - The name of the endpoint.

* `backend_port` - The port number on the compute node.

* `protocol` - The protocol of the endpoint.

* `frontend_port_range` - The range of external ports that are used to provide inbound access to the backendPort on the individual compute nodes in the format of `1000-1100`.

* `network_security_group_rules` - The list of network security group rules that are applied to the endpoint.

---

A `network_security_group_rules` block exports the following:

* `access` - The action that should be taken for a specified IP address, subnet range or tag.

* `priority` - The priority for this rule.

* `source_address_prefix` - The source address prefix or tag to match for the rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Pool.
