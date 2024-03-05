---
subcategory: "Nginx"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_nginx_deployment"
description: |-
  Gets information about an existing Nginx Deployment.
---

# Data Source: azurerm_nginx_deployment

Use this data source to access information about an existing Nginx Deployment.

## Example Usage

```hcl
data "azurerm_nginx_deployment" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_nginx_deployment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Nginx Deployment.

* `resource_group_name` - (Required) The name of the Resource Group where the Nginx Deployment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Nginx Deployment.

* `capacity` - The number of NGINX capacity units for this Nginx Deployment.

* `diagnose_support_enabled` - Whether diagnostic settings are enabled.

* `email` - Preferred email associated with the Nginx Deployment.

* `frontend_private` - A `frontend_private` block as defined below.

* `frontend_public` - A `frontend_public` block as defined below.

* `identity` - A `identity` block as defined below.

* `ip_address` - The IP address of the Nginx Deployment.

* `location` - The Azure Region where the Nginx Deployment exists.

* `logging_storage_account` - A `logging_storage_account` block as defined below.

* `managed_resource_group` - Auto-generated managed resource group for the Nginx Deployment.

* `network_interface` - A `network_interface` block as defined below.

* `nginx_version` - NGINX version of the Nginx Deployment.

* `sku` - Name of the SKU for this Nginx Deployment.

* `automatic_upgrade_channel` - The automatic upgrade channel for this NGINX deployment.

* `tags` - A mapping of tags assigned to the Nginx Deployment.

---

A `frontend_private` block exports the following:

* `allocation_method` - The method of allocating the private IP to the Nginx Deployment.

* `ip_address` - Private IP address of the Nginx Deployment.

* `subnet_id` - The subnet ID of the Nginx Deployment.

---

A `frontend_public` block exports the following:

* `ip_address` - List of public IPs of the Ngix Deployment.

---

A `identity` block exports the following:

* `identity_ids` - List of identities attached to the Nginx Deployment.

* `type` - Type of identity attached to the Nginx Deployment.

---

A `logging_storage_account` block exports the following:

* `container_name` - the container name of Storage Account for logging.

* `name` - The account name of the StorageAccount for logging.

---

A `network_interface` block exports the following:

* `subnet_id` - The subnet resource ID of the Nginx Deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Nginx Deployment.
