---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
description: |-
  Manages an Azure Batch pool.

---

# azurerm_batch_pool

Manages an Azure Batch pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "testaccbatch"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "testaccsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "example" {
  name                                = "testaccbatch"
  resource_group_name                 = azurerm_resource_group.example.name
  location                            = azurerm_resource_group.example.location
  pool_allocation_mode                = "BatchService"
  storage_account_id                  = azurerm_storage_account.example.id
  storage_account_authentication_mode = "StorageKeys"

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_certificate" "example" {
  resource_group_name  = azurerm_resource_group.example.name
  account_name         = azurerm_batch_account.example.name
  certificate          = filebase64("certificate.cer")
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "example" {
  name                = "testaccpool"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_batch_account.example.name
  display_name        = "Test Acc Pool Auto"
  vm_size             = "STANDARD_A1_V2"
  node_agent_sku_id   = "batch.node.ubuntu 20.04"

  auto_scale {
    evaluation_interval = "PT15M"

    formula = <<EOF
      startingNumberOfVMs = 1;
      maxNumberofVMs = 25;
      pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
      pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 *   TimeInterval_Second));
      $TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
EOF

  }

  storage_image_reference {
    publisher = "microsoft-azure-batch"
    offer     = "ubuntu-server-container"
    sku       = "20-04-lts"
    version   = "latest"
  }

  container_configuration {
    type = "DockerCompatible"
    container_registries {
      registry_server = "docker.io"
      user_name       = "login"
      password        = "apassword"
    }
  }

  start_task {
    command_line       = "echo 'Hello World from $env'"
    task_retry_maximum = 1
    wait_for_success   = true

    common_environment_properties = {
      env = "TEST"
    }

    user_identity {
      auto_user {
        elevation_level = "NonAdmin"
        scope           = "Task"
      }
    }
  }

  certificate {
    id             = azurerm_batch_certificate.example.id
    store_location = "CurrentUser"
    visibility     = ["StartTask"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch pool. Changing this forces a new resource to be created.

* `account_name` - (Required) Specifies the name of the Batch account in which the pool will be created. Changing this forces a new resource to be created.

* `node_agent_sku_id` - (Required) Specifies the SKU of the node agents that will be created in the Batch pool. Changing this forces a new resource to be created.

* `stop_pending_resize_operation` - (Optional) Whether to stop if there is a pending resize operation on this pool.

* `vm_size` - (Required) Specifies the size of the VM created in the Batch pool. Changing this forces a new resource to be created.

* `storage_image_reference` - (Required) A `storage_image_reference` block for the virtual machines that will compose the Batch pool as defined below. Changing this forces a new resource to be created.

* `data_disks` - (Optional) A `data_disks` block describes the data disk settings as defined below.

* `display_name` - (Optional) Specifies the display name of the Batch pool. Changing this forces a new resource to be created.

* `disk_encryption` - (Optional) A `disk_encryption` block, as defined below, describes the disk encryption configuration applied on compute nodes in the pool. Disk encryption configuration is not supported on Linux pool created with Virtual Machine Image or Shared Image Gallery Image.

* `extensions` - (Optional) An `extensions` block as defined below.

* `inter_node_communication` - (Optional) Whether the pool permits direct communication between nodes. This imposes restrictions on which nodes can be assigned to the pool. Enabling this value can reduce the chance of the requested number of nodes to be allocated in the pool. Values allowed are `Disabled` and `Enabled`. Defaults to `Enabled`.

* `identity` - (Optional) An `identity` block as defined below.

* `license_type` - (Optional) The type of on-premises license to be used when deploying the operating system. This only applies to images that contain the Windows operating system, and should only be used when you hold valid on-premises licenses for the nodes which will be deployed. If omitted, no on-premises licensing discount is applied. Values are: "Windows_Server" - The on-premises license is for Windows Server. "Windows_Client" - The on-premises license is for Windows Client.

* `max_tasks_per_node` - (Optional) Specifies the maximum number of tasks that can run concurrently on a single compute node in the pool. Defaults to `1`. Changing this forces a new resource to be created.

* `fixed_scale` - (Optional) A `fixed_scale` block that describes the scale settings when using fixed scale as defined below.

* `auto_scale` - (Optional) A `auto_scale` block that describes the scale settings when using auto scale as defined below.

* `start_task` - (Optional) A `start_task` block that describes the start task settings for the Batch pool as defined below.

* `certificate` - (Optional) One or more `certificate` blocks that describe the certificates to be installed on each compute node in the pool as defined below.

* `container_configuration` - (Optional) The container configuration used in the pool's VMs. One `container_configuration` block as defined below.

* `metadata` - (Optional) A map of custom batch pool metadata.

* `mount` - (Optional) A `mount` block defined as below.

* `network_configuration` - (Optional) A `network_configuration` block that describes the network configurations for the Batch pool as defined below. Changing this forces a new resource to be created.

* `node_placement` - (Optional) A `node_placement` block that describes the placement policy for allocating nodes in the pool as defined below.

* `os_disk_placement` - (Optional) Specifies the ephemeral disk placement for operating system disk for all VMs in the pool. This property can be used by user in the request to choose which location the operating system should be in. e.g., cache disk space for Ephemeral OS disk provisioning. For more information on Ephemeral OS disk size requirements, please refer to Ephemeral OS disk size requirements for Windows VMs at <https://docs.microsoft.com/en-us/azure/virtual-machines/windows/ephemeral-os-disks#size-requirements> and Linux VMs at <https://docs.microsoft.com/en-us/azure/virtual-machines/linux/ephemeral-os-disks#size-requirements>. The only possible value is `CacheDisk`.

* `security_profile` - (Optional) A `security_profile` block that describes the security settings for the Batch pool as defined below. Changing this forces a new resource to be created.

* `target_node_communication_mode` - (Optional) The desired node communication mode for the pool. Possible values are `Classic`, `Default` and `Simplified`.

* `task_scheduling_policy` - (Optional) A `task_scheduling_policy` block that describes how tasks are distributed across compute nodes in a pool as defined below. If not specified, the default is spread as defined below.

* `user_accounts` - (Optional) A `user_accounts` block that describes the list of user accounts to be created on each node in the pool as defined below.

* `windows` - (Optional) A `windows` block that describes the Windows configuration in the pool as defined below.

-> **Note:** For Windows compute nodes, the Batch service installs the certificates to the specified certificate store and location. For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of `remoteUser`, a `certs` directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

~> **Note:** `fixed_scale` and `auto_scale` blocks cannot be used both at the same time.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Batch Account. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Batch Account.

---

A `data_disks` block supports the following:

* `lun` - (Required) The lun is used to uniquely identify each data disk. If attaching multiple disks, each should have a distinct lun. The value must be between 0 and 63, inclusive.

* `caching` - (Optional) Values are: "none" - The caching mode for the disk is not enabled. "readOnly" - The caching mode for the disk is read only. "readWrite" - The caching mode for the disk is read and write. For information about the caching options see: <https://blogs.msdn.microsoft.com/windowsazurestorage/2012/06/27/exploring-windows-azure-drives-disks-and-images/>. Possible values are `None`, `ReadOnly` and `ReadWrite`. Defaults to `ReadOnly`.

* `disk_size_gb` - (Required) The initial disk size in GB when creating new data disk.

* `storage_account_type` - (Optional) The storage account type to be used for the data disk. Values are: Possible values are `Standard_LRS` - The data disk should use standard locally redundant storage. `Premium_LRS` - The data disk should use premium locally redundant storage. Defaults to `Standard_LRS`.

---

A `disk_encryption` block supports the following:

The disk encryption configuration applied on compute nodes in the pool. Disk encryption configuration is not supported on Linux pool created with Virtual Machine Image or Shared Image Gallery Image.

* `disk_encryption_target` - (Required) On Linux pool, only \"TemporaryDisk\" is supported; on Windows pool, \"OsDisk\" and \"TemporaryDisk\" must be specified.

---

An `extensions` block supports the following:

The virtual machine extension for the pool.
If specified, the extensions mentioned in this configuration will be installed on each node.

* `name` - (Required) The name of the virtual machine extension.

* `publisher` - (Required) The name of the extension handler publisher.The name of the extension handler publisher.

* `type` - (Required) The type of the extensions.

* `type_handler_version` - (Optional) The version of script handler.

* `auto_upgrade_minor_version` - (Optional) Indicates whether the extension should use a newer minor version if one is available at deployment time. Once deployed, however, the extension will not upgrade minor versions unless redeployed, even with this property set to true.

* `automatic_upgrade_enabled` - (Optional) Indicates whether the extension should be automatically upgraded by the platform if there is a newer version available. Supported values are `true` and `false`.

~> **Note:** When `automatic_upgrade_enabled` is set to `true`, the `type_handler_version` is automatically updated by the Azure platform when a new version is available and any change in `type_handler_version` should be manually ignored by user.

* `settings_json` - (Optional) JSON formatted public settings for the extension, the value should be encoded with [`jsonencode`](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function.

* `protected_settings` - (Optional) JSON formatted protected settings for the extension, the value should be encoded with [`jsonencode`](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function. The extension can contain either `protected_settings` or `provision_after_extensions` or no protected settings at all.

* `provision_after_extensions` - (Optional) The collection of extension names. Collection of extension names after which this extension needs to be provisioned.

---

A `node_placement` block supports the following:

Node placement Policy type on Batch Pools. Allocation policy used by Batch Service to provision the nodes. If not specified, Batch will use the regional policy.

* `policy` - (Optional) The placement policy for allocating nodes in the pool. Values are: "Regional": All nodes in the pool will be allocated in the same region; "Zonal": Nodes in the pool will be spread across different zones with the best effort balancing. Defaults to `Regional`.

---

A `storage_image_reference` block supports the following:

This block provisions virtual machines in the Batch Pool from one of two sources: an Azure Platform Image (e.g. Ubuntu/Windows Server) or a Custom Image.

To provision from an Azure Platform Image, the following fields are applicable:

* `publisher` - (Optional) Specifies the publisher of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `offer` - (Optional) Specifies the offer of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `sku` - (Optional) Specifies the SKU of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `version` - (Optional) Specifies the version of the image used to create the virtual machines. Changing this forces a new resource to be created.

To provision a Custom Image, the following fields are applicable:

* `id` - (Optional) Specifies the ID of the Custom Image which the virtual machines should be created from. Changing this forces a new resource to be created. See [official documentation](https://docs.microsoft.com/azure/batch/batch-custom-images) for more details.

---

A `fixed_scale` block supports the following:

* `node_deallocation_method` - (Optional) It determines what to do with a node and its running task(s) if the pool size is decreasing. Values are `Requeue`, `RetainedData`, `TaskCompletion` and `Terminate`.

* `target_dedicated_nodes` - (Optional) The number of nodes in the Batch pool. Defaults to `1`.

* `target_low_priority_nodes` - (Optional) The number of low priority nodes in the Batch pool. Defaults to `0`.

* `resize_timeout` - (Optional) The timeout for resize operations. Defaults to `PT15M`.

---

A `auto_scale` block supports the following:

* `evaluation_interval` - (Optional) The interval to wait before evaluating if the pool needs to be scaled. Defaults to `PT15M`.

* `formula` - (Required) The autoscale formula that needs to be used for scaling the Batch pool.

---

A `start_task` block supports the following:

* `command_line` - (Required) The command line executed by the start task.

* `container` - (Optional) A `container` block is the settings for the container under which the start task runs as defined below. When this is specified, all directories recursively below the `AZ_BATCH_NODE_ROOT_DIR` (the root of Azure Batch directories on the node) are mapped into the container, all task environment variables are mapped into the container, and the task command line is executed in the container.

* `task_retry_maximum` - (Optional) The number of retry count.

* `wait_for_success` - (Optional) A flag that indicates if the Batch pool should wait for the start task to be completed. Default to `false`.

* `common_environment_properties` - (Optional) A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - (Required) A `user_identity` block that describes the user identity under which the start task runs as defined below.

* `resource_file` - (Optional) One or more `resource_file` blocks that describe the files to be downloaded to a compute node as defined below.

---

A `container` block supports the following:

* `image_name` - (Required) The image to use to create the container in which the task will run. This is the full image reference, as would be specified to "docker pull". If no tag is provided as part of the image name, the tag ":latest" is used as a default.

* `run_options` - (Optional) Additional options to the container create command. These additional options are supplied as arguments to the "docker create" command, in addition to those controlled by the Batch Service.

* `registry` - (Optional) The `container_registries` block defined as below.

* `working_directory` - (Optional) A flag to indicate where the container task working directory is. Possible values are `TaskWorkingDirectory` and `ContainerImageDefault`.

---

A `user_identity` block supports the following:

* `user_name` - (Optional) The username to be used by the Batch pool start task.

* `auto_user` - (Optional) A `auto_user` block that describes the user identity under which the start task runs as defined below.

~> **Note:** `user_name` and `auto_user` blocks cannot be used both at the same time, but you need to define one or the other.

---

A `auto_user` block supports the following:

* `elevation_level` - (Optional) The elevation level of the user identity under which the start task runs. Possible values are `Admin` or `NonAdmin`. Defaults to `NonAdmin`.

* `scope` - (Optional) The scope of the user identity under which the start task runs. Possible values are `Task` or `Pool`. Defaults to `Task`.

---

A `certificate` block supports the following:

* `id` - (Required) The ID of the Batch Certificate to install on the Batch Pool, which must be inside the same Batch Account.

* `store_location` - (Required) The location of the certificate store on the compute node into which to install the certificate. Possible values are `CurrentUser` or `LocalMachine`.

-> **Note:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of `remoteUser`, a 'certs' directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

* `store_name` - (Optional) The name of the certificate store on the compute node into which to install the certificate. This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). Common store names include: `My`, `Root`, `CA`, `Trust`, `Disallowed`, `TrustedPeople`, `TrustedPublisher`, `AuthRoot`, `AddressBook`, but any custom store name can also be used.

* `visibility` - (Optional) Which user accounts on the compute node should have access to the private data of the certificate. Possible values are `StartTask`, `Task` and `RemoteUser`.

---

A `container_configuration` block supports the following:

* `type` - (Optional) The type of container configuration. Possible value is `DockerCompatible`.

* `container_image_names` - (Optional) A list of container image names to use, as would be specified by `docker pull`. Changing this forces a new resource to be created.

* `container_registries` - (Optional) One or more `container_registries` blocks as defined below. Additional container registries from which container images can be pulled by the pool's VMs. Changing this forces a new resource to be created.

---

A `resource_file` block supports the following:

* `auto_storage_container_name` - (Optional) The storage container name in the auto storage account.

* `blob_prefix` - (Optional) The blob prefix to use when downloading blobs from an Azure Storage container. Only the blobs whose names begin with the specified prefix will be downloaded. The property is valid only when `auto_storage_container_name` or `storage_container_url` is used. This prefix can be a partial filename or a subdirectory. If a prefix is not specified, all the files in the container will be downloaded.

* `file_mode` - (Optional) The file permission mode represented as a string in octal format (e.g. `"0644"`). This property applies only to files being downloaded to Linux compute nodes. It will be ignored if it is specified for a `resource_file` which will be downloaded to a Windows node. If this property is not specified for a Linux node, then a default value of 0770 is applied to the file.

* `file_path` - (Optional) The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified, `file_path` is optional and is the directory to download the files to. In the case where `file_path` is used as a directory, any directory structure already associated with the input data will be retained in full and appended to the specified filePath directory. The specified relative path cannot break out of the task's working directory (for example by using '..').

* `http_url` - (Optional) The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read permissions on the blob, or set the ACL for the blob or its container to allow public access.

* `storage_container_url` - (Optional) The URL of the blob container within Azure Blob Storage. This URL must be readable and listable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read and list permissions on the blob, or set the ACL for the blob or its container to allow public access.

* `user_assigned_identity_id` - (Optional) An identity reference from pool's user assigned managed identity list.

~> **Note:** Exactly one of `auto_storage_container_name`, `storage_container_url` and `auto_user` must be specified.

---

A `container_registries` block supports the following:

* `registry_server` - (Required) The container registry URL. Changing this forces a new resource to be created.

* `user_name` - (Optional) The user name to log into the registry server. Changing this forces a new resource to be created.

* `password` - (Optional) The password to log into the registry server. Changing this forces a new resource to be created.

* `user_assigned_identity_id` - (Optional) The reference to the user assigned identity to use to access an Azure Container Registry instead of username and password. Changing this forces a new resource to be created.

---

An `mount` block supports the following:

Any property below is mutually exclusive with all other properties.

* `azure_blob_file_system` - (Optional) A `azure_blob_file_system` block defined as below.

* `azure_file_share` - (Optional) A `azure_file_share` block defined as below.

* `cifs_mount` - (Optional) A `cifs_mount` block defined as below.

* `nfs_mount` - (Optional) A `nfs_mount` block defined as below.

---

An `azure_blob_file_system` block supports the following:

* `account_name` - (Required) The Azure Storage Account name.

* `container_name` - (Required) The Azure Blob Storage Container name.

* `relative_mount_path` - (Required) The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `account_key` - (Optional) The Azure Storage Account key. This property is mutually exclusive with both `sas_key` and `identity_id`; exactly one must be specified.

* `sas_key` - (Optional) The Azure Storage SAS token. This property is mutually exclusive with both `account_key` and `identity_id`; exactly one must be specified.

* `identity_id` - (Optional) The ARM resource id of the user assigned identity. This property is mutually exclusive with both `account_key` and `sas_key`; exactly one must be specified.

* `blobfuse_options` - (Optional) Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

An `azure_file_share` block supports the following:

* `account_name` - (Required) The Azure Storage Account name.

* `account_key` - (Required) The Azure Storage Account key.

* `azure_file_url` - (Required) The Azure Files URL. This is of the form 'https://{account}.file.core.windows.net/'.

* `relative_mount_path` - (Required) The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - (Optional) Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `cifs_mount` block supports the following:

* `user_name` - (Required) The user to use for authentication against the CIFS file system.

* `password` - (Required) The password to use for authentication against the CIFS file system.

* `source` - (Required) The URI of the file system to mount.

* `relative_mount_path` - (Required) The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - (Optional) Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `nfs_mount` block supports the following:

* `source` - (Required) The URI of the file system to mount.

* `relative_mount_path` - (Required) The relative path on compute node where the file system will be mounted All file systems are mounted relative to the Batch mounts directory, accessible via the `AZ_BATCH_NODE_MOUNTS_DIR` environment variable.

* `mount_options` - (Optional) Additional command line options to pass to the mount command. These are 'net use' options in Windows and 'mount' options in Linux.

---

A `network_configuration` block supports the following:

* `subnet_id` - (Optional) The ARM resource identifier of the virtual network subnet which the compute nodes of the pool will join. Changing this forces a new resource to be created.

* `dynamic_vnet_assignment_scope` - (Optional) The scope of dynamic vnet assignment. Allowed values: `none`, `job`. Changing this forces a new resource to be created. Defaults to `none`.

* `accelerated_networking_enabled` - (Optional) Whether to enable accelerated networking. Possible values are `true` and `false`. Defaults to `false`. Changing this forces a new resource to be created.

* `public_ips` - (Optional) A list of public IP ids that will be allocated to nodes. Changing this forces a new resource to be created.

* `endpoint_configuration` - (Optional) A list of `endpoint_configuration` blocks that can be used to address specific ports on an individual compute node externally as defined below. Set as documented in the inbound_nat_pools block below. Changing this forces a new resource to be created.

* `public_address_provisioning_type` - (Optional) Type of public IP address provisioning. Supported values are `BatchManaged`, `UserManaged` and `NoPublicIPAddresses`.

---

A `endpoint_configuration` block supports the following:

* `name` - (Required) The name of the endpoint. The name must be unique within a Batch pool, can contain letters, numbers, underscores, periods, and hyphens. Names must start with a letter or number, must end with a letter, number, or underscore, and cannot exceed 77 characters. Changing this forces a new resource to be created.

* `backend_port` - (Required) The port number on the compute node. Acceptable values are between `1` and `65535` except for `29876`, `29877` as these are reserved. Changing this forces a new resource to be created.

* `protocol` - (Required) The protocol of the endpoint. Acceptable values are `TCP` and `UDP`. Changing this forces a new resource to be created.

* `frontend_port_range` - (Required) The range of external ports that will be used to provide inbound access to the backendPort on individual compute nodes in the format of `1000-1100`. Acceptable values range between `1` and `65534` except ports from `50000` to `55000` which are reserved by the Batch service. All ranges within a pool must be distinct and cannot overlap. Values must be a range of at least `100` nodes. Changing this forces a new resource to be created.

* `network_security_group_rules` - (Optional) A list of `network_security_group_rules` blocks as defined below that will be applied to the endpoint. The maximum number of rules that can be specified across all the endpoints on a Batch pool is `25`. If no network security group rules are specified, a default rule will be created to allow inbound access to the specified backendPort. Set as documented in the network_security_group_rules block below. Changing this forces a new resource to be created.

---

A `network_security_group_rules` block supports the following:

* `access` - (Required) The action that should be taken for a specified IP address, subnet range or tag. Acceptable values are `Allow` and `Deny`. Changing this forces a new resource to be created.

* `priority` - (Required) The priority for this rule. The value must be at least `150`. Changing this forces a new resource to be created.

* `source_address_prefix` - (Required) The source address prefix or tag to match for the rule. Changing this forces a new resource to be created.

* `source_port_ranges` - (Optional) The source port ranges to match for the rule. Valid values are `*` (for all ports 0 - 65535) or arrays of ports or port ranges (i.e. `100-200`). The ports should in the range of 0 to 65535 and the port ranges or ports can't overlap. If any other values are provided the request fails with HTTP status code 400. Default value will be `*`. Changing this forces a new resource to be created.

---

A `task_scheduling_policy` block supports the following:

* `node_fill_type` - (Optional) Supported values are "Pack" and "Spread". "Pack" means as many tasks as possible (taskSlotsPerNode) should be assigned to each node in the pool before any tasks are assigned to the next node in the pool. "Spread" means that tasks should be assigned evenly across all nodes in the pool.

---
A `security_profile` block supports the following:

* `host_encryption_enabled` - (Optional) Whether to enable host encryption for the Virtual Machine or Virtual Machine Scale Set. This will enable the encryption for all the disks including Resource/Temp disk at host itself. Possible values are `true` and `false`. Changing this forces a new resource to be created.

* `security_type` - (Optional) The security type of the Virtual Machine. Possible values are `confidentialVM` and `trustedLaunch`. Changing this forces a new resource to be created.

* `secure_boot_enabled` - (Optional) Whether to enable secure boot for the Virtual Machine or Virtual Machine Scale Set. Possible values are `true` and `false`. Changing this forces a new resource to be created.

* `vtpm_enabled` - (Optional) Whether to enable virtual trusted platform module (vTPM) for the Virtual Machine or Virtual Machine Scale Set. Possible values are `true` and `false`. Changing this forces a new resource to be created.

~> **Note:** `security_profile` block can only be specified during creation and does not support updates.

~> **Note:** `security_type` must be specified to set UEFI related properties including `secure_boot_enabled` and `vtpm_enabled`.

---

A `user_accounts` block supports the following:

* `name` - (Required) The name of the user account.

* `password` - (Required) The password for the user account.

* `elevation_level` - (Required) The elevation level of the user account. "NonAdmin" - The auto user is a standard user without elevated access. "Admin" - The auto user is a user with elevated access and operates with full Administrator permissions. The default value is nonAdmin.

* `linux_user_configuration` - (Optional) The `linux_user_configuration` block defined below is a linux-specific user configuration for the user account. This property is ignored if specified on a Windows pool. If not specified, the user is created with the default options.

* `windows_user_configuration` - (Optional) The `windows_user_configuration` block defined below is a windows-specific user configuration for the user account. This property can only be specified if the user is on a Windows pool. If not specified and on a Windows pool, the user is created with the default options.

---

A `linux_user_configuration` block supports the following:

* `uid` - (Optional) The group ID for the user account. The `uid` and `gid` properties must be specified together or not at all. If not specified the underlying operating system picks the gid.

* `gid` - (Optional) The user ID of the user account. The `uid` and `gid` properties must be specified together or not at all. If not specified the underlying operating system picks the uid.

* `ssh_private_key` - (Optional) The SSH private key for the user account. The private key must not be password protected. The private key is used to automatically configure asymmetric-key based authentication for SSH between nodes in a Linux pool when the pool's enableInterNodeCommunication property is true (it is ignored if enableInterNodeCommunication is false). It does this by placing the key pair into the user's .ssh directory. If not specified, password-less SSH is not configured between nodes (no modification of the user's .ssh directory is done).

---

A `windows_user_configuration` block supports the following:

* `login_mode` - (Required) Specifies login mode for the user. The default value for VirtualMachineConfiguration pools is interactive mode and for CloudServiceConfiguration pools is batch mode. Values supported are "Batch" and "Interactive".

---

A `windows` block supports the following:

Windows operating system settings on the virtual machine. This property must not be specified if the imageReference specifies a Linux OS image.

* `enable_automatic_updates` - (Optional) Whether automatic updates are enabled on the virtual machine. Defaults to `true`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Batch Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Batch Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Batch Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Batch Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Batch Pool.

## Import

Batch Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_batch_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Batch/batchAccounts/myBatchAccount1/pools/myBatchPool1
```
