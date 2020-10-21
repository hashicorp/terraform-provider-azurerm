---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management"
description: |-
  Gets information about an existing API Management Service.
---

# Data Source: azurerm_api_management

Use this data source to access information about an existing API Management Service.

## Example Usage

```hcl
data "azurerm_api_management" "example" {
  name                = "search-api"
  resource_group_name = "search-service"
}

output "api_management_id" {
  value = data.azurerm_api_management.example.id
}
```

## Argument Reference

* `name` - The name of the API Management service.

* `resource_group_name` - The Name of the Resource Group in which the API Management Service exists.

## Attributes Reference

* `id` - The ID of the API Management Service.

* `additional_location` - Zero or more `additional_location` blocks as defined below

* `location` - The Azure location where the API Management Service exists.

* `gateway_url` - The URL for the API Management Service's Gateway.

* `gateway_regional_url` - The URL for the Gateway in the Default Region.

* `identity` - (Optional) An `identity` block as defined below.

* `hostname_configuration` - A `hostname_configuration` block as defined below.

* `management_api_url` - The URL for the Management API.

* `notification_sender_email` - The email address from which the notification will be sent.

* `portal_url` - The URL of the Publisher Portal.

* `developer_portal_url` - The URL for the Developer Portal associated with this API Management service.

* `public_ip_addresses` - The Public IP addresses of the API Management Service.

* `private_ip_addresses` - The Private IP addresses of the API Management Service.

* `publisher_name` - The name of the Publisher/Company of the API Management Service.

* `publisher_email` - The email of Publisher/Company of the API Management Service.

* `scm_url` - The SCM (Source Code Management) endpoint.

* `sku` - A `sku` block as documented below.

* `tags` - A mapping of tags assigned to the resource.

---

A `additional_location` block exports the following:

* `location` - The location name of the additional region among Azure Data center regions.

* `gateway_regional_url` - Gateway URL of the API Management service in the Region.

* `public_ip_addresses` - Public Static Load Balanced IP addresses of the API Management service in the additional location. Available only for Basic, Standard and Premium SKU.

* `private_ip_addresses` - Private IP addresses of the API Management service in the additional location, for instances using virtual network mode.

---

A `identity` block exports the following:

~> **Note:** User Assigned Managed Identities are in Preview

* `type` - Specifies the type of Managed Service Identity that is configured on this API Management Service.

* `principal_id` - Specifies the Principal ID of the System Assigned Managed Service Identity that is configured on this API Management Service.

* `tenant_id` - Specifies the Tenant ID of the System Assigned Managed Service Identity that is configured on this API Management Service.

* `identity_ids` - A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `hostname_configuration` block exports the following:

* `management` - One or more `management` blocks as documented below.

* `portal` - One or more `portal` blocks as documented below.

* `developer_portal` - One or more `developer_portal` blocks as documented below.

* `proxy` - One or more `proxy` blocks as documented below.

* `scm` - One or more `scm` blocks as documented below.

---

A `management` block exports the following:

* `host_name` - The Hostname used for the Management API.

* `key_vault_id` - The ID of the Key Vault Secret which contains the SSL Certificate.

* `negotiate_client_certificate` - Is Client Certificate Negotiation enabled?

---

A `portal` block exports the following:

* `host_name` - The Hostname used for the Portal.

* `key_vault_id` - The ID of the Key Vault Secret which contains the SSL Certificate.

* `negotiate_client_certificate` - Is Client Certificate Negotiation enabled?

---

A `developer_portal` block exports the following:

* `host_name` - The Hostname used for the Portal.

* `key_vault_id` - The ID of the Key Vault Secret which contains the SSL Certificate.

* `negotiate_client_certificate` - Is Client Certificate Negotiation enabled?

---

A `proxy` block exports the following:

* `default_ssl_binding` - Is this the default SSL Binding?

* `host_name` - The Hostname used for the Proxy.

* `key_vault_id` - The ID of the Key Vault Secret which contains the SSL Certificate.

* `negotiate_client_certificate` - Is Client Certificate Negotiation enabled?

---

A `scm` block exports the following:

* `host_name` - The Hostname used for the SCM URL.

* `key_vault_id` - The ID of the Key Vault Secret which contains the SSL Certificate.

* `negotiate_client_certificate` - Is Client Certificate Negotiation enabled?


---

A `sku` block exports the following:

* `name` - Specifies the plan's pricing tier.

* `capacity` - Specifies the number of units associated with this API Management service.

---


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Service.
