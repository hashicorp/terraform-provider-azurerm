---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
description: |-
  Manage an Azure Spring Cloud Application.
---

# azurerm_spring_cloud_app

Manage an Azure Spring Cloud Application.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

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

* `ingress_settings` - (Optional) An `ingress_settings` block as defined below.

* `persistent_disk` - (Optional) An `persistent_disk` block as defined below.

* `public_endpoint_enabled` - (Optional) Should the App in vnet injection instance exposes endpoint which could be accessed from Internet?

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

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ingress_settings` block supports the following:

* `backend_protocol` - (Optional) Specifies how ingress should communicate with this app backend service. Allowed values are `GRPC` and `Default`. Defaults to `Default`.

* `read_timeout_in_seconds` - (Optional) Specifies the ingress read time out in seconds. Defaults to `300`.

* `send_timeout_in_seconds` - (Optional) Specifies the ingress send time out in seconds. Defaults to `60`.

* `session_affinity` - (Optional) Specifies the type of the affinity, set this to `Cookie` to enable session affinity. Allowed values are `Cookie` and `None`. Defaults to `None`.

* `session_cookie_max_age` - (Optional) Specifies the time in seconds until the cookie expires.

---

An `persistent_disk` block supports the following:

* `size_in_gb` - (Required) Specifies the size of the persistent disk in GB. Possible values are between `0` and `50`.

* `mount_path` - (Optional) Specifies the mount path of the persistent disk. Defaults to `/persistent`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application.

* `fqdn` - The Fully Qualified DNS Name of the Spring Application in the service.

* `url` - The public endpoint of the Spring Cloud Application.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application.

## Import

Spring Cloud Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.AppPlatform/spring/myservice/apps/myapp
```
