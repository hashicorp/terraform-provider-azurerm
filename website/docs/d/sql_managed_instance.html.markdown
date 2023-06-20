---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_managed_instance"
description: |-
  Gets information about a SQL Managed Instance.
---

# Data Source: azurerm_sql_managed_instance

Use this data source to access information about an existing SQL Managed Instance.

-> **Note:** The `azurerm_sql_managed_instance` data source is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_managed_instance`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/mssql_managed_instance) data source instead.

## Example Usage

```hcl
data "azurerm_sql_managed_instance" "example" {
  name                = "example_mi"
  resource_group_name = "example-resources"
}

output "sql_instance_id" {
  value = data.azurerm_sql_managed_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance.

* `resource_group_name` - (Required) The name of the Resource Group in which the SQL Managed Instance exists.

## Attributes Reference

* `id` - The SQL Managed Instance ID.

* `fqdn` - The fully qualified domain name of the Azure Managed SQL Instance.

* `location` - Location where the resource exists.

* `sku_name` - SKU Name for the SQL Managed Instance.

* `vcores` - Number of cores assigned to your instance.

* `storage_size_in_gb` - Maximum storage space for your instance.

* `license_type` - Type of license the Managed Instance uses.

* `administrator_login` - The administrator login name for the new server.

* `subnet_id` - The subnet resource id that the SQL Managed Instance is associated with.

* `collation` - Specifies how the SQL Managed Instance is collated.

* `public_data_endpoint_enabled` - Is the public data endpoint enabled?

* `minimum_tls_version` - The Minimum TLS Version.

* `proxy_override` - How the SQL Managed Instance is accessed.

* `timezone_id` - The TimeZone ID that the SQL Managed Instance is operating in.

* `dns_zone_partner_id` - The ID of the Managed Instance which is sharing the DNS zone.

* `identity` - An `identity` block as defined below.

* `storage_account_type` - Storage account type used to store backups for this SQL Managed Instance.

* `tags` - A mapping of tags assigned to the resource.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Managed Instance.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Managed Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Azure Managed Instance.
