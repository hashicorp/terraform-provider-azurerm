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

The following attributes are exported:

* `id` - The Batch pool ID.

* `name` - The name of the Batch pool.

* `account_name` - The name of the Batch account.

* `vm_size` - The size of the VM created in the Batch pool.

* `auto_scale` - A `auto_scale` block that describes the scale settings when using auto scale.

* `certificate` - One or more `certificate` blocks that describe the certificates installed on each compute node in the pool.

* `container_configuration` - The container configuration used in the pool's VMs.

* `data_disks` - A `data_disks` block describes the data disk settings.

* `display_name` - Specifies the display name of the Batch pool.

* `disk_encryption_configuration` - A `disk_encryption_configuration` block describes the disk encryption configuration applied on compute nodes in the pool.

* `extensions` - An `extensions` block describes the extension settings.

* `fixed_scale` - A `fixed_scale` block that describes the scale settings when using fixed scale.

* `identity` - An `identity` block describes the identity settings.

* `inter_node_communication` - Whether the pool permits direct communication between nodes. This imposes restrictions on which nodes can be assigned to the pool. Enabling this value can reduce the chance of the requested number of nodes to be allocated in the pool. If not specified, this value defaults to "Disabled". Value is "Disabled" or "Enabled".

* `license_type` - The type of on-premises license to be used when deploying the operating system. This only applies to images that contain the Windows operating system, and should only be used when you hold valid on-premises licenses for the nodes which will be deployed. If omitted, no on-premises licensing discount is applied. Values are: Windows_Server - The on-premises license is for Windows Server. Windows_Client - The on-premises license is for Windows Client.

* `max_tasks_per_node` - Specifies the maximum number of tasks that can run concurrently on a single compute node in the pool. Defaults to `1`. Changing this forces a new resource to be created.

* `mount_configuration` - A `mount_configuration` block that describes mount configuration.

* `network_configuration` - A `mount_configuration` block that describes network configuration.

* `node_agent_sku_id` - The SKU of the node agents in the Batch pool.

* `node_placement_configuration` - A `node_placement_configuration` block that describes the placement policy for allocating nodes in the pool.

* `os_disk_placement_setting` - Specifies the ephemeral disk placement for operating system disk for all VMs in the pool. This property can be used by user in the request to choose which location the operating system should be in. e.g., cache disk space for Ephemeral OS disk provisioning. For more information on Ephemeral OS disk size requirements, please refer to Ephemeral OS disk size requirements for Windows VMs at https://docs.microsoft.com/en-us/azure/virtual-machines/windows/ephemeral-os-disks#size-requirements and Linux VMs at https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ephemeral-os-disks#size-requirements.

* `start_task` - A `start_task` block that describes the start task settings for the Batch pool.

* `storage_image_reference` - The reference of the storage image used by the nodes in the Batch pool.

* `task_scheduling_policy` - A `task_scheduling_policy` block that describes how tasks are distributed across compute nodes in a pool.

* `user_accounts` - A `user_accounts` block that describes the list of user accounts to be created on each node in the pool.

* `windows_configuration` - A `windows_configuration` block that describes the Windows configuration in the pool.

---

A `auto_scale` block exports the following:

* `evaluation_interval` - The interval to wait before evaluating if the pool needs to be scaled.

* `formula` - The autoscale formula that needs to be used for scaling the Batch pool.

---

A `certificate` block exports the following:

* `id` - The fully qualified ID of the certificate installed on the pool.

* `store_location` - The location of the certificate store on the compute node into which the certificate is installed, either `CurrentUser` or `LocalMachine`.

-> **NOTE:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of 'remoteUser', a 'certs' directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

* `store_name` - The name of the certificate store on the compute node into which the certificate is installed.

-> **NOTE:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference).

* `visibility` - Which user accounts on the compute node have access to the private data of the certificate.

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

* `identity_id` - The ARM resource id of the user assigned identity.

---

A `data_disks` block exports the following:

* `lun` - The lun is used to uniquely identify each data disk. If attaching multiple disks, each should have a distinct lun. The value must be between 0 and 63, inclusive.

* `caching` - Values are: "none" - The caching mode for the disk is not enabled. "readOnly" - The caching mode for the disk is read only. "readWrite" - The caching mode for the disk is read and write. The default value for caching is "none". For information about the caching options see: https://blogs.msdn.microsoft.com/windowsazurestorage/2012/06/27/exploring-windows-azure-drives-disks-and-images/.

* `disk_size_gb` - The initial disk size in GB when creating new data disk.

* `storage_account_type` - The storage account type to be used for the data disk. If omitted, the default is "Standard_LRS". Values are: "Standard_LRS" - The data disk should use standard locally redundant storage. "Premium_LRS" - The data disk should use premium locally redundant storage.

---

A `disk_encryption_configuration` block exports the following:

The disk encryption configuration applied on compute nodes in the pool. Disk encryption configuration is not supported on Linux pool created with Virtual Machine Image or Shared Image Gallery Image.

* `disk_encryption_target` - On Linux pool, only \"TemporaryDisk\" is supported; on Windows pool, \"OsDisk\" and \"TemporaryDisk\" must be specified.

---

An `extensions` block exports the following:

The virtual machine extension for the pool.
If specified, the extensions mentioned in this configuration will be installed on each node.

* `name` - The name of the virtual machine extension.

* `publisher` - The name of the extension handler publisher.The name of the extension handler publisher.

* `type` - The type of the extensions.

* `type_handler_version` - The version of script handler.

* `auto_upgrade_minor_version` - Indicates whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.

* `settings` - JSON formatted public settings for the extension.

* `protected_settings` - The extension can contain either `protected_settings` or `provision_after_extensions` or no protected settings at all.

* `provision_after_extensions` - The collection of extension names. Collection of extension names after which this extension needs to be provisioned.

---

A `fixed_scale` block exports the following:

* `node_deallocation_option` - Determines what to do with a node and its running task(s) after it has been selected for deallocation.

* `target_dedicated_nodes` - The number of nodes in the Batch pool.

* `target_low_priority_nodes` - The number of low priority nodes in the Batch pool.

* `resize_timeout` - The timeout for resize operations.

---

A `node_placement_configuration` block exports the following:

Node placement Policy type on Batch Pools. Allocation policy used by Batch Service to provision the nodes. If not specified, Batch will use the regional policy.

* `policy` - The placement policy for allocating nodes in the pool. Values are: "Regional": All nodes in the pool will be allocated in the same region; "Zonal": Nodes in the pool will be spread across different zones with the best effort balancing.

---

An `mount_configuration` exports the following:

Any property below is mutually exclusive with all other properties.

* `azure_blob_file_system_configuration` - A `azure_blob_file_system_configuration` block defined as below.

* `azure_file_share_configuration` - A `azure_file_share_configuration` block defined as below.

* `cifs_mount_configuration` - A `cifs_mount_configuration` block defined as below.

* `nfs_mount_configuration` - A `nfs_mount_configuration` block defined as below.

---

An `azure_blob_file_system_configuration` block exports the following:

* `account_name` - The Azure Storage Account name.

* `container_name` - The Azure Blob Storage Container name.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `account_key` - The Azure Storage Account key. This property is mutually exclusive with both `sas_key` and `identity_id`; exactly one must be specified.

* `sas_key` - The Azure Storage SAS token. This property is mutually exclusive with both `account_key` and `identity_id`; exactly one must be specified.

* `identity_id` - The ARM resource id of the user assigned identity. This property is mutually exclusive with both `account_key` and `sas_key`; exactly one must be specified.

* `blobfuse_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

An `azure_file_share_configuration` block exports the following:

* `account_name` - The Azure Storage Account name.

* `account_key` - The Azure Storage Account key.

* `azure_file_url` - The Azure Files URL. This is of the form 'https://{account}.file.core.windows.net/'.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `cifs_mount_configuration` block exports the following:

* `user_name` - The user to use for authentication against the CIFS file system.

* `password` - The password to use for authentication against the CIFS file system.

* `source` - The URI of the file system to mount.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `nfs_mount_configuration` block exports the following:

* `source` - The URI of the file system to mount.

* `relative_mount_path` - The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `network_configuration` block exports the following:

* `subnet_id` - The ARM resource identifier of the virtual network subnet which the compute nodes of the pool will join. Changing this forces a new resource to be created.

* `dynamic_vnet_assignment_scope` - The scope of dynamic vnet assignment. Allowed values: `None`, `Job`.

* `public_ips` - A list of public IP ids that will be allocated to nodes. Changing this forces a new resource to be created.

* `endpoint_configuration` - A list of inbound NAT pools that can be used to address specific ports on an individual compute node externally. Set as documented in the inbound_nat_pools block below. Changing this forces a new resource to be created.

* `public_address_provisioning_type` - Type of public IP address provisioning. Supported values are `BatchManaged`, `UserManaged` and `NoPublicIPAddresses`.

---

A `endpoint_configuration` block exports the following:

* `name` - The name of the endpoint. The name must be unique within a Batch pool, can contain letters, numbers, underscores, periods, and hyphens. Names must start with a letter or number, must end with a letter, number, or underscore, and cannot exceed 77 characters. Changing this forces a new resource to be created.

* `backend_port` - The port number on the compute node. Acceptable values are between `1` and `65535` except for `29876`, `29877` as these are reserved. Changing this forces a new resource to be created.

* `protocol` - The protocol of the endpoint. Acceptable values are `TCP` and `UDP`. Changing this forces a new resource to be created.

* `frontend_port_range` - The range of external ports that will be used to provide inbound access to the backendPort on individual compute nodes in the format of `1000-1100`. Acceptable values range between `1` and `65534` except ports from `50000` to `55000` which are reserved by the Batch service. All ranges within a pool must be distinct and cannot overlap. Values must be a range of at least `100` nodes. Changing this forces a new resource to be created.

* `network_security_group_rules` - A list of network security group rules that will be applied to the endpoint. The maximum number of rules that can be specified across all the endpoints on a Batch pool is `25`. If no network security group rules are specified, a default rule will be created to allow inbound access to the specified backendPort. Set as documented in the network_security_group_rules block below. Changing this forces a new resource to be created.

---

A `network_security_group_rules` block exports the following:

* `access` - The action that should be taken for a specified IP address, subnet range or tag. Acceptable values are `Allow` and `Deny`. Changing this forces a new resource to be created.

* `priority` - The priority for this rule. The value must be at least `150`. Changing this forces a new resource to be created.

* `source_address_prefix` - The source address prefix or tag to match for the rule. Changing this forces a new resource to be created.

* `source_port_ranges` - The source port ranges to match for the rule. Valid values are '*' (for all ports 0 - 65535) or arrays of ports or port ranges (i.e. 100-200). The ports should in the range of 0 to 65535 and the port ranges or ports can't overlap. If any other values are provided the request fails with HTTP status code 400. Default value will be *.

---

A `node_placement_configuration` block exports the following:

* `policy` - The placement policy for allocating nodes in the pool.

---

A `start_task` block exports the following:

* `command_line` - The command line executed by the start task.

* `container_settings` - The settings for the container under which the start task runs.
  
* `task_retry_maximum` - The number of retry count.

* `wait_for_success` - A flag that indicates if the Batch pool should wait for the start task to be completed.

* `common_environment_properties` - A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - A `user_identity` block that describes the user identity under which the start task runs.

* `resource_file` - One or more `resource_file` blocks that describe the files to be downloaded to a compute node.

---

A `container_settings` block exports the following:

* `image_name` - The image to use to create the container in which the task will run. 

* `container_run_options` - Additional options to the container create command. 

* `registry` - The same reference as `container_registry` block defined as follows.

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

A `resource_file` block exports the following:

* `auto_storage_container_name` - The storage container name in the auto storage account.

* `blob_prefix` - The blob prefix used when downloading blobs from an Azure Storage container.

* `file_mode` - The file permission mode attribute represented as a string in octal format (e.g. `"0644"`).

* `file_path` - The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified.

* `http_url` - The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access.

* `storage_container_url` - The URL of the blob container within Azure Blob Storage.

* `identity_id` - The reference to the user assigned identity to use to access Azure Blob Storage.

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

---

A `task_scheduling_policy` block exports the following:

* `node_fill_type` - Supported values are "Pack" and "Spread". "Pack" means as many tasks as possible (taskSlotsPerNode) should be assigned to each node in the pool before any tasks are assigned to the next node in the pool. "Spread" means that tasks should be assigned evenly across all nodes in the pool.

---

A `user_accounts` block supports the following:

* `name` - The name of the user account.

* `password` - The password for the user account.

* `elevation_level` - The elevation level of the user account. "NonAdmin" - The auto user is a standard user without elevated access. "Admin" - The auto user is a user with elevated access and operates with full Administrator permissions. The default value is nonAdmin.

* `linux_user_configuration` - The `linux_user_configuration` block defined below is a linux-specific user configuration for the user account. This property is ignored if specified on a Windows pool. If not specified, the user is created with the default options.

* `windows_user_configuration` - The `windows_user_configuration` block defined below is a windows-specific user configuration for the user account. This property can only be specified if the user is on a Windows pool. If not specified and on a Windows pool, the user is created with the default options.

---

A `linux_user_configuration` block supports the following:

* `uid` - The group ID for the user account. 

* `gid` - The user ID of the user account. 

* `ssh_private_key` - The SSH private key for the user account. 

---

A `windows_user_configuration` block supports the following:

* `login_mode` - Specifies login mode for the user. 

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Pool.
