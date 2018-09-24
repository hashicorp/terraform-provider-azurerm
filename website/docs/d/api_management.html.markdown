---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
sidebar_current: "docs-azurerm-datasource-azurerm_api_management"
description: |-
  Get information about an API Management service.
---

# Data Source: azurerm_api_management

Use this data source to obtain information about an API Management service.

## Example Usage

```hcl
data "azurerm_api_management" "test" {
  name                = "search-api"
  resource_group_name = "search-service"
}

output "api_management_id" {
  value = "${data.azurerm_api_management.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the API Management service.

* `resource_group_name` - (Required) The Name of the Resource Group where the API Management service exists.

## Attributes Reference

* `id` - The ID of the API Management Service.

* `location` - The Azure location where the API Management Service exists.

* `publisher_name` - The name of publisher/company.

* `publisher_email` - The email of publisher/company.

* `sku` - A `sku` block as documented below.

* `notification_sender_email` - Email address from which the notification will be sent.

* `gateway_url` - Gateway URL of the API Management service.

* `gateway_regional_url` - Gateway URL of the API Management service in the Default Region.

* `portal_url` - Publisher portal endpoint Url of the API Management service.

* `management_api_url` - Management API endpoint URL of the API Management service.

* `scm_url` - SCM endpoint URL of the API Management service.

* `additional_location` - Additional datacenter locations of the API Management service. The `additional_location` block is documented below.

* `hostname_configurations` - Custom hostname configuration of the API Management service. The `hostname_configurations` block is documented below.

* `tags` - A mapping of tags assigned to the resource.

---

A `sku` block supports the following:

* `name` - Specifies the plan's pricing tier.

* `capacity` - Specifies the number of units associated with this API Management service.


A `additional_location` block supports the following:

* `location` - The location name of the additional region among Azure Data center regions.

* `gateway_regional_url` - Gateway URL of the API Management service in the Region.

* `static_ips` - Static IP addresses of the location's virtual machines.

`hostname_configurations` block supports the following:

* `management` - The `management` block is documented below.

* `portal` - The `portal` block is documented below.

* `proxy` - The `proxy` block is documented below.

* `scm` - The `scm` block is documented below.

`management`, `portal` and `scm` blocks supports the following:

* `host_name` - Hostname to configure on the Api Management service.

* `key_vault_id` - Url to the KeyVault Secret containing the SSL Certificate. If absolute Url containing version is provided, auto-update of ssl certificate will not work. This requires the `identity` attribute to be set. The secret should be of type `application/x-pkcs12`.
