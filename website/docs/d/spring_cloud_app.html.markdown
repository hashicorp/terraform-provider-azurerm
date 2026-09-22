---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app"
description: |-
  Gets information about an existing Spring Cloud Application
---

# Data Source: azurerm_spring_cloud_app

Use this data source to access information about an existing Spring Cloud Application.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app` data source is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
data "azurerm_spring_cloud_app" "example" {
  name                = azurerm_spring_cloud_app.example.name
  resource_group_name = azurerm_spring_cloud_app.example.resource_group_name
  service_name        = azurerm_spring_cloud_app.example.service_name
}

output "spring_cloud_app_id" {
  value = data.azurerm_spring_cloud_app.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Application.

* `resource_group_name` - (Required) The name of the Resource Group where the Spring Cloud Application exists.

* `service_name` - (Required) The name of the Spring Cloud Service.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Spring Cloud Application.

* `addon_json` - A JSON object that contains the addon configurations of the Spring Cloud Service.

* `custom_persistent_disk` - A `custom_persistent_disk` block as defined below.

* `fqdn` - The Fully Qualified DNS Name.

* `https_only` - Is only HTTPS allowed?

* `identity` - An `identity` block as defined below.

* `ingress_settings` - An `ingress_settings` block as defined below.

* `is_public` - Does the Spring Cloud Application have public endpoint?

* `persistent_disk` - A `persistent_disk` block as defined below.

* `public_endpoint_enabled` - Does the App in vnet injection instance expose an endpoint which could be accessed from the Internet?

* `url` - The public endpoint of the Spring Cloud Application.

* `tls_enabled` - Is End to End TLS Enabled?

---

The `custom_persistent_disk` block exports the following:

* `storage_name` - The name of the Spring Cloud Storage.

* `mount_path` - The mount path of the persistent disk.

* `share_name` - The share name of the Azure File share.

* `mount_options` - The mount options for the persistent disk.

* `read_only_enabled` - Whether the persistent disk is a readOnly one.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Spring Cloud Application.

* `type` - The Type of Managed Identity assigned to the Spring Cloud Application.

---

The `ingress_settings` block exports the following:

* `backend_protocol` - How ingress communicates with this app backend service.

* `read_timeout_in_seconds` - The ingress read time out in seconds.

* `send_timeout_in_seconds` - The ingress send time out in seconds.

* `session_affinity` - The type of the affinity.

* `session_cookie_max_age` - The time in seconds until the cookie expires.

---

The `persistent_disk` block exports the following:

* `mount_path` - The mount path of the persistent disk.

* `size_in_gb` - The size of the persistent disk in GB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application.
