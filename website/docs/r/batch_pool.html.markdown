---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_pool"
sidebar_current: "docs-azurerm-resource-batch-pool"
description: |-
  Manages an Azure Batch pool.

---

# azurerm_batch_pool

Manages an Azure Batch pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "testaccsa"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags = {
    env = "test"
  }
}

resource "azurerm_batch_certificate" "testcer" {
  resource_group_name  = "${azurerm_resource_group.test.name}"
  account_name         = "${azurerm_batch_account.test.name}"
  certificate          = "${filebase64("certificate.cer")}"
  format               = "Cer"
  thumbprint           = "312d31a79fa0cef49c00f769afc2b73e9f4edf34"
  thumbprint_algorithm = "SHA1"
}

resource "azurerm_batch_pool" "test" {
  name                = "testaccpool"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_batch_account.test.name}"
  display_name        = "Test Acc Pool Auto"
  vm_size             = "Standard_A1"
  node_agent_sku_id   = "batch.node.ubuntu 16.04"

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
    sku       = "16-04-lts"
    version   = "latest"
  }

  container_configuration {
    type = "DockerCompatible"
    container_registries = [
      {
        registry_server = "docker.io"
        user_name       = "login"
        password        = "apassword"
      },
    ]
  }

  start_task {
    command_line         = "echo 'Hello World from $env'"
    max_task_retry_count = 1
    wait_for_success     = true

    environment = {
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
    id         = "${azurerm_batch_certificate.testcer.id}"
    visibility = ["StartTask"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Batch pool. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Batch pool. Changing this forces a new resource to be created.

~> **NOTE:** To work around [a bug in the Azure API](https://github.com/Azure/azure-rest-api-specs/issues/5574) this property is currently treated as case-insensitive. A future version of Terraform will require that the casing is correct.

* `account_name` - (Required) Specifies the name of the Batch account in which the pool will be created. Changing this forces a new resource to be created.

* `node_agent_sku_id` - (Required) Specifies the Sku of the node agents that will be created in the Batch pool.

* `vm_size` - (Required) Specifies the size of the VM created in the Batch pool.

* `storage_image_reference` - (Required) A `storage_image_reference` for the virtual machines that will compose the Batch pool.

* `display_name` - (Optional) Specifies the display name of the Batch pool.

* `max_tasks_per_node` - (Optional) Specifies the maximum number of tasks that can run concurrently on a single compute node in the pool. Defaults to `1`. Changing this forces a new resource to be created.

* `fixed_scale` - (Optional) A `fixed_scale` block that describes the scale settings when using fixed scale.

* `auto_scale` - (Optional) A `auto_scale` block that describes the scale settings when using auto scale.

* `start_task` - (Optional) A `start_task` block that describes the start task settings for the Batch pool.

* `certificate` - (Optional) One or more `certificate` blocks that describe the certificates to be installed on each compute node in the pool.

* `container_configuration` - (Optional) The container configuration used in the pool's VMs.

-> **NOTE:** For Windows compute nodes, the Batch service installs the certificates to the specified certificate store and location. For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of `remoteUser`, a `certs` directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

~> **Please Note:** `fixed_scale` and `auto_scale` blocks cannot be used both at the same time.

---
A `storage_image_reference` block supports the following:

This block provisions virtual machines in the Batch Pool from one of two sources: an Azure Platform Image (e.g. Ubuntu/Windows Server) or a Custom Image.

To provision from an Azure Platform Image, the following fields are applicable:

* `publisher` - (Required) Specifies the publisher of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `offer` - (Required) Specifies the offer of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `sku` - (Required) Specifies the SKU of the image used to create the virtual machines. Changing this forces a new resource to be created.

* `version` - (Optional) Specifies the version of the image used to create the virtual machines. Changing this forces a new resource to be created.

To provision a Custom Image, the following fields are applicable:

* `id` - (Required) Specifies the ID of the Custom Image which the virtual machines should be created from. Changing this forces a new resource to be created. See [official documentation](https://docs.microsoft.com/en-us/azure/batch/batch-custom-images) for more details.
---

A `fixed_scale` block supports the following:

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

* `max_task_retry_count` - (Optional) The number of retry count. Defaults to `1`.

* `wait_for_success` - (Optional) A flag that indicates if the Batch pool should wait for the start task to be completed. Default to `false`.

* `environment` - (Optional) A map of strings (key,value) that represents the environment variables to set in the start task.

* `user_identity` - (Required) A `user_identity` block that describes the user identity under which the start task runs.

* `resource_file` - (Optional) One or more `resource_file` blocks that describe the files to be downloaded to a compute node.

---

A `user_identity` block supports the following:

* `user_name` - (Optional) The username to be used by the Batch pool start task.

* `auto_user` - (Optional) A `auto_user` block that describes the user identity under which the start task runs.

~> **Please Note:** `user_name` and `auto_user` blocks cannot be used both at the same time, but you need to define one or the other.

---

A `auto_user` block supports the following:

* `elevation_level` - (Optional) The elevation level of the user identity under which the start task runs. Possible values are `Admin` or `NonAdmin`. Defaults to `NonAdmin`.

* `scope` - (Optional) The scope of the user identity under which the start task runs. Possible values are `Task` or `Pool`. Defaults to `Task`.

---

A `certificate` block supports the following:

* `id` - (Required) The ID of the Batch Certificate to install on the Batch Pool, which must be inside the same Batch Account.

* `store_location` - (Required) The location of the certificate store on the compute node into which to install the certificate. Possible values are `CurrentUser` or `LocalMachine`.

 -> **NOTE:** This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). For Linux compute nodes, the certificates are stored in a directory inside the task working directory and an environment variable `AZ_BATCH_CERTIFICATES_DIR` is supplied to the task to query for this location. For certificates with visibility of `remoteUser`, a 'certs' directory is created in the user's home directory (e.g., `/home/{user-name}/certs`) and certificates are placed in that directory.

* `store_name` - (Optional) The name of the certificate store on the compute node into which to install the certificate. This property is applicable only for pools configured with Windows nodes (that is, created with cloudServiceConfiguration, or with virtualMachineConfiguration using a Windows image reference). Common store names include: `My`, `Root`, `CA`, `Trust`, `Disallowed`, `TrustedPeople`, `TrustedPublisher`, `AuthRoot`, `AddressBook`, but any custom store name can also be used. The default value is `My`.

* `visibility` - (Optional) Which user accounts on the compute node should have access to the private data of the certificate.

---

A `container_configuration` block supports the following:

* `type` - (Optional) The type of container configuration. Possible value is `DockerCompatible`.

* `container_registries` - (Optional) Additional container registries from which container images can be pulled by the pool's VMs.

---

A `resource_file` block supports the following:

* `auto_storage_container_name` - (Optional) The storage container name in the auto storage account.

* `blob_prefix` - (Optional) The blob prefix to use when downloading blobs from an Azure Storage container. Only the blobs whose names begin with the specified prefix will be downloaded. The property is valid only when `auto_storage_container_name` or `storage_container_url` is used. This prefix can be a partial filename or a subdirectory. If a prefix is not specified, all the files in the container will be downloaded.

* `file_mode` - (Optional) The file permission mode represented as a string in octal format (e.g. `"0644"`). This property applies only to files being downloaded to Linux compute nodes. It will be ignored if it is specified for a `resource_file` which will be downloaded to a Windows node. If this property is not specified for a Linux node, then a default value of 0770 is applied to the file.

* `file_path` - (Optional) The location on the compute node to which to download the file, relative to the task's working directory. If the `http_url` property is specified, the `file_path` is required and describes the path which the file will be downloaded to, including the filename. Otherwise, if the `auto_storage_container_name` or `storage_container_url` property is specified, `file_path` is optional and is the directory to download the files to. In the case where `file_path` is used as a directory, any directory structure already associated with the input data will be retained in full and appended to the specified filePath directory. The specified relative path cannot break out of the task's working directory (for example by using '..').

* `http_url` - (Optional) The URL of the file to download. If the URL is Azure Blob Storage, it must be readable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read permissions on the blob, or set the ACL for the blob or its container to allow public access.

* `storage_container_url` - (Optional) The URL of the blob container within Azure Blob Storage. This URL must be readable and listable using anonymous access; that is, the Batch service does not present any credentials when downloading the blob. There are two ways to get such a URL for a blob in Azure storage: include a Shared Access Signature (SAS) granting read and list permissions on the blob, or set the ACL for the blob or its container to allow public access.

~> **Please Note:** Exactly one of `auto_storage_container_name`, `storage_container_url` and `auto_user` must be specified.

---

A `container_registries` block supports the following:

* `registry_server` - (Optional) The container registry URL. The default is "docker.io". Changing this forces a new resource to be created.

* `user_name` - (Optional) The user name to log into the registry server. Changing this forces a new resource to be created.

* `password` - (Optional) The password to log into the registry server. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The Batch pool ID.
