---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mssql_managed_instance"
description: |-
  Gets information about a Microsoft SQL Azure Managed Instance.
---

# Data Source: azurerm_mssql_managed_instance

Use this data source to access information about an existing Microsoft SQL Azure Managed Instance.

## Example Usage

```hcl
data "azurerm_mssql_managed_instance" "example" {
  name                = "managedsqlinstance"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SQL Managed Instance.

* `resource_group_name` - (Required) The name of the resource group where the SQL Managed Instance exists.

## Attributes Reference

The following attributes are exported:

* `administrator_login` - The administrator login name for the SQL Managed Instance.

* `collation` - Specifies how the SQL Managed Instance will be collated.

* `customer_managed_key` - Specifies KeyVault key, used by SQL Managed Instance for Transparent Data Encryption.

* `dns_zone` - The Dns Zone where the SQL Managed Instance is located.

* `dns_zone_partner_id` - The ID of the SQL Managed Instance which shares the DNS zone.

* `fqdn` - The fully qualified domain name of the Azure Managed SQL Instance.

* `id` - The SQL Managed Instance ID.

* `identity` - An `identity` block as defined below.

* `license_type` - What type of license the SQL Managed Instance uses.

* `location` - Specifies the supported Azure location where the resource exists.

* `minimum_tls_version` - The Minimum TLS Version.

* `proxy_override` - Specifies how the SQL Managed Instance will be accessed.

* `public_data_endpoint_enabled` - Whether the public data endpoint is enabled.

* `sku_name` - Specifies the SKU Name of the SQL Managed Instance.

* `storage_account_type` - Specifies the storage account type used to store backups for this database.

* `storage_size_in_gb` - Maximum storage space allocated for the SQL Managed Instance.

* `subnet_id` - The subnet resource ID that the SQL Managed Instance is associated with.

* `tags` - A mapping of tags assigned to the resource.

* `timezone_id` - The TimeZone ID that the SQL Managed Instance is running in.

* `vcores` - Number of cores that are assigned to the SQL Managed Instance.

---

The `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Managed Instance.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Managed Instance.

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned with the Identity of this SQL Managed Instance.

* `type` - The identity type of the SQL Managed Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Managed Instance.
