---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
description: |-
  Manage an Azure Spring Cloud Application.
---

# azurerm_spring_cloud_app

Manage an Azure Spring Cloud Application.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example-springcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name

  identity {
    type = "SystemAssigned"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group in which to create the Spring Cloud Application. Changing this forces a new resource to be created.

* `service_name` - (Required) Specifies the name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `addon_json` - (Optional) A JSON object that contains the addon configurations of the Spring Cloud Service.

* `custom_persistent_disk` - (Optional) A `custom_persistent_disk` block as defined below.
  
* `identity` - (Optional) An `identity` block as defined below.

* `is_public` - (Optional) Does the Spring Cloud Application have public endpoint? Defaults to `false`.

* `https_only` - (Optional) Is only HTTPS allowed? Defaults to `false`.

* `persistent_disk` - (Optional) An `persistent_disk` block as defined below.

* `tls_enabled` - (Optional) Is End to End TLS Enabled? Defaults to `false`.

---
An `custom_persistent_disk` block exports the following:

* `storage_name` - (Required) The name of the Spring Cloud Storage.

* `mount_path` - (Required) The mount path of the persistent disk.

* `share_name` - (Required) The share name of the Azure File share.

* `mount_options` - (Optional) These are the mount options for a persistent disk.

* `read_only_enabled` - (Optional) Indicates whether the persistent disk is a readOnly one.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Spring Cloud Application. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Spring Cloud Application.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `persistent_disk` block supports the following:

* `size_in_gb` - (Required) Specifies the size of the persistent disk in GB. Possible values are between `0` and `50`.

* `mount_path` - (Optional) Specifies the mount path of the persistent disk. Defaults to `/persistent`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Application.

* `fqdn` - The Fully Qualified DNS Name of the Spring Application in the service.

* `url` - The public endpoint of the Spring Cloud Application.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application.

## Import

Spring Cloud Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/Spring/myservice/apps/myapp
```
