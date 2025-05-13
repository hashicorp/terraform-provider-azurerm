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

* `node_agent_sku_id` - The SKU of the node agents in the Batch pool.

* `vm_size` - The size of the VM created in the Batch pool.

* `fixed_scale` - A `fixed_scale` block that describes the scale settings when using fixed scale.

* `auto_scale` - A `auto_scale` block that describes the scale settings when using auto scale.

* `data_disks` - A `data_disks` block describes the data disk settings.

* `disk_encryption` - A `disk_encryption` block describes the disk encryption configuration applied on compute nodes in the pool.

* `extensions` - An `extensions` block describes the extension settings

* `inter_node_communication` - Whether the pool permits direct communication between nodes. This imposes restrictions on which nodes can be assigned to the pool. Enabling this value can reduce the chance of the requested number of nodes to be allocated in the pool.

* `license_type` - The type of on-premises license to be used when deploying the operating system.

* `node_placement` - A `node_placement` block that describes the placement policy for allocating nodes in the pool.

* `os_disk_placement` - Specifies the ephemeral disk placement for operating system disk for all VMs in the pool.

* `storage_image_reference` - The reference of the storage image used by the nodes in the Batch pool.

* `start_task` - A `start_task` block that describes the start task settings for the Batch pool.

* `task_scheduling_policy` - A `task_scheduling_policy` block that describes how tasks are distributed across compute nodes in a pool.

* `user_accounts` - A `user_accounts` block that describes the list of user accounts to be created on each node in the pool.

* `windows` - A `windows` block that describes the Windows configuration in the pool.

* `max_tasks_per_node` - The maximum number of tasks that can run concurrently on a single compute node in the pool.

* `mount` - A `mount` block that describes mount configuration.

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

A `data_disks` block exports the following:

* `lun` - The lun is used to uniquely identify each data disk.

* `caching` - The caching mode of data disks.

* `disk_size_gb` - The initial disk size in GB when creating new data disk.

* `storage_account_type` - The storage account type to be used for the data disk.

---

A `disk_encryption` block exports the following:

* `disk_encryption_target` - On Linux pool, only `TemporaryDisk` is supported; on Windows pool, `OsDisk` and `TemporaryDisk` must be specified.

---

An `extensions` block exports the following:

* `name` - The name of the virtual machine extension.

* `publisher` - The name of the extension handler publisher.The name of the extension handler publisher.

* `type` - The type of the extensions.

* `type_handler_version` - The version of script handler.

* `auto_upgrade_minor_version` - Indicates whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.

* `settings_json` - JSON formatted public settings for the extension.

* `protected_settings` - The extension can contain either `protected_settings` or `provision_after_extensions` or no protected settings at all.

* `provision_after_extensions` - The collection of extension names. Collection of extension names after which this extension needs to be provisioned.

---

A `node_placement` block exports the following:

* `policy` - The placement policy for allocating nodes in the pool.

---

A `start_task` block exports the following:

* `command_line` - The command line executed by the start task.

* `container` - The settings for the container under which the start task runs.

* `task_retry_maximum` - The number of retry count

* `wait_for_success` - A flag that indicates if the Batch pool should wait for the start task to be completed.

* `common_environment_properties` - A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - A `user_identity` block that describes the user identity under which the start task runs.

* `resource_file` - One or more `resource_file` blocks that describe the files to be downloaded to a compute node.

---

A `container` block exports the following:

* `image_name` - The image to use to create the container in which the task will run.

* `run_options` - Additional options to the container create command.

* `registry` - The same reference as `container_registries` block defined as follows.

* `working_directory` - A flag to indicate where the container task working directory is.

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

-> **Note:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of 'remoteUser', a 'certs' directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

* `store_name` - The name of the certificate store on the compute node into which the certificate is installed.

-> **Note:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference).

* `visibility` - Which user accounts on the compute node have access to the private data of the certificate.

---

A `resource_file` block exports the following:

* `auto_storage_container_name` - The storage container name in the auto storage account.

* `blob_prefix` - The blob prefix used when downloading blobs from an Azure Storage container.

* `file_mode` - The file permission mode attribute represented as a string in octal format (e.g. `"0644"`).

* `file_path` - The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified.

* `http_url` - The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access.

* `storage_container_url` - The URL of the blob container within Azure Blob Storage.

* `user_assigned_identity_id` - An identity reference from pool's user assigned managed identity list.

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

* `user_assigned_identity_id` - The reference to the user assigned identity to use to access an Azure Container Registry instead of username and password.

---

An `mount` exports the following:

Any property below is mutually exclusive with all other properties.

* `azure_blob_file_system` - A `azure_blob_file_system` block defined as below.

* `azure_file_share` - A `azure_file_share` block defined as below.

* `cifs_mount` - A `cifs_mount` block defined as below.

* `nfs_mount` - A `nfs_mount` block defined as below.

---

An `azure_blob_file_system` block exports the following:

* `account_name` - The Azure Storage Account name.

* `container_name` - The Azure Blob Storage Container name.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `account_key` - The Azure Storage Account key. This property is mutually exclusive with both `sas_key` and `identity_id`; exactly one must be specified.

* `sas_key` - The Azure Storage SAS token. This property is mutually exclusive with both `account_key` and `identity_id`; exactly one must be specified.

* `identity_id` - The ARM resource id of the user assigned identity. This property is mutually exclusive with both `account_key` and `sas_key`; exactly one must be specified.

* `blobfuse_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

An `azure_file_share` block exports the following:

* `account_name` - The Azure Storage Account name.

* `account_key` - The Azure Storage Account key.

* `azure_file_url` - The Azure Files URL. This is of the form 'https://{account}.file.core.windows.net/'.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `cifs_mount` block exports the following:

* `user_name` - The user to use for authentication against the CIFS file system.

* `password` - The password to use for authentication against the CIFS file system.

* `source` - The URI of the file system to mount.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `nfs_mount` block exports the following:

* `source` - The URI of the file system to mount.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `network_configuration` block exports the following:

* `subnet_id` - The ARM resource identifier of the virtual network subnet which the compute nodes of the pool are joined too.

* `dynamic_vnet_assignment_scope` - The scope of dynamic vnet assignment.

* `endpoint_configuration` - The inbound NAT pools that are used to address specific ports on the individual compute node externally.

* `public_ips` - A list of public IP ids that will be allocated to nodes.

* `public_address_provisioning_type` - Type of public IP address provisioning.

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

* `source_port_ranges` - The source port ranges to match for the rule.

---

A `task_scheduling_policy` block exports the following:

* `node_fill_type` - Supported values are `Pack` and `Spread`. `Pack` means as many tasks as possible (taskSlotsPerNode) should be assigned to each node in the pool before any tasks are assigned to the next node in the pool. `Spread` means that tasks should be assigned evenly across all nodes in the pool.

---

A `user_accounts` block exports the following:

* `name` - The name of the user account.

* `password` - The password for the user account.

* `elevation_level` - The elevation level of the user account. "NonAdmin" - The auto user is a standard user without elevated access. "Admin" - The auto user is a user with elevated access and operates with full Administrator permissions. The default value is nonAdmin.

* `linux_user_configuration` - The `linux_user_configuration` block defined below is a linux-specific user configuration for the user account. This property is ignored if specified on a Windows pool. If not specified, the user is created with the default options.

* `windows_user_configuration` - The `windows_user_configuration` block defined below is a windows-specific user configuration for the user account. This property can only be specified if the user is on a Windows pool. If not specified and on a Windows pool, the user is created with the default options.

---

A `linux_user_configuration` block exports the following:

* `uid` - The group ID for the user account.

* `gid` - The user ID of the user account.

* `ssh_private_key` - The SSH private key for the user account.

---

A `windows_user_configuration` block exports the following:

* `login_mode` - Specifies login mode for the user.

---

A `windows` block exports the following:

Windows operating system settings on the virtual machine. This property must not be specified if the imageReference specifies a Linux OS image.

* `enable_automatic_updates` - Whether automatic updates are enabled on the virtual machine.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Pool.
