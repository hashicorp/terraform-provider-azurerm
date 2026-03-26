---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_environment_storage"
description: |-
  Gets information about a Container App Environment Storage.
---

# Data Source: azurerm_container_app_environment_storage

Use this data source to access information about an existing Container App Environment Storage.

## Example Usage

```hcl
data "azurerm_container_app_environment" "example" {
  name                = "existing-environment"
  resource_group_name = "existing-resources"
}

data "azurerm_container_app_environment_storage" "example" {
  name                         = "existing-storage"
  container_app_environment_id = data.azurerm_container_app_environment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Container App Environment Storage.

* `container_app_environment_id` - (Required) The ID of the Container App Environment to which this storage belongs.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Container App Environment Storage.

* `access_mode` - The access mode to connect this storage to the Container App.

* `account_name` - The Azure Storage Account in which the Share is located.

* `nfs_server_url` - The NFS server URL for the Azure File Share.

* `share_name` - The name of the Azure Storage Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Container App Environment Storage.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.App` - 2025-07-01
