---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance"
description: |-
  Gets information about an existing SQL elastic pool.
---

# Data Source: azurerm_mssql_managed_instance

Use this data source to access information about an existing managed instance.

## Example Usage

```hcl
data "azurerm_mssql_managed_instance" "example" {
  name                = "managedinstancename"
  resource_group_name = "example-resources"
}

output "managedInstance" {
  value = data.azurerm_mssql_managed_instance.example
}
```

## Argument Reference

* `name` - The name of the managed instance.

* `resource_group_name` - The name of the resource group which contains the managed instance.


## Attributes Reference

* `administrator_login` - The managed instance admin login username.

* `location` - Specifies the supported Azure location where the resource exists.

* `collation` - The managed instance SQL collation.
 
* `dns_zone` - The DNS Zone to which this managed instance belongs.

* `instance_pool_id` - The resource id of the instance pool to which this managed instance belongs to.

* `license_type` - The license type of the managed instance.

* `maintenance_configuration_id` - The resource id of the maintenance configuration for the managed instance.

* `minimal_tls_version` - Whether or not this elastic pool is zone redundant.

* `proxy_override` - Whether or not this elastic pool is zone redundant.

* `public_data_endpoint_enabled` - Whether or not this elastic pool is zone redundant.

* `storage_size_gb` - Whether or not this elastic pool is zone redundant.

* `subnet_id` - Whether or not this elastic pool is zone redundant.

* `timezone_id` - Whether or not this elastic pool is zone redundant.

* `vcores` - Whether or not this elastic pool is zone redundant.

* `fully_qualified_domain_name` - Whether or not this elastic pool is zone redundant.

* `type` - Whether or not this elastic pool is zone redundant.

* `tags` - Whether or not this elastic pool is zone redundant.

* `identity` - An `identity` block as defined below.

---

An `identity` block supports the following:

* `type` -  The identity type of the Microsoft SQL Managed instance.

* `principal_id` - The identity principal id of the managed instance.

* `tenant_id` - The identity tenant id of the managed instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the managed instance.
