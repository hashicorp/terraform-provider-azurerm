---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
sidebar_current: "docs-azurerm-datasource-batch-pool"
description: |-
  Get information about an existing Azure Batch pool.

---

# Data source: azurerm_batch_pool

Use this data source to access information about an existing Batch pool

## Example Usage

```hcl
data "azurerm_batch_pool "test" {
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

* `resource_file` - (Optional) One or more `resource_file` blocks that describe the files to be downloaded to a compute node. 

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

* `blob_prefix` - The blob prefix used when downloading blobs from an Azure Storage container. Only the blobs whose names begin with the specified prefix will be downloaded. If a prefix is not specified, all the files in the container will be downloaded.

* `file_Mode` - The file permission mode attribute in octal format. This property applies only to files being downloaded to Linux compute nodes. It will be ignored if it is specified for a `resource_file` which will be downloaded to a Windows node. If this property is not specified for a Linux node, then a default value of 0770 is applied to the file.

* `file_path` - The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified, `file_path` is optional and is the directory to download the files to. In the case where `file_path` is used as a directory, any directory structure already associated with the input data will be retained in full and appended to the specified filePath directory. The specified relative path cannot break out of the task's working directory (for example by using '..').

* `http_url` - The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read permissions on the blob, or set the ACL for the blob or its container to allow public access.

* `storage_container_url` - The URL of the blob container within Azure Blob Storage. This URL must be readable and listable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read and list permissions on the blob, or set the ACL for the blob or its container to allow public access.
