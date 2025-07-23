---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_soa_record"
description: |-
  Gets information about an existing Private DNS SOA Record.
---

# Data Source: azurerm_private_dns_soa_record

Use this data source to access information about an existing Private DNS SOA Record within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

## Example Usage

```hcl
resource "azurerm_private_dns_soa_record" "example" {
  zone_name           = "test-zone"
  resource_group_name = "test-rg"
}

output "private_dns_soa_record_id" {
  value = data.azurerm_private_dns_soa_record.example.id
}
```

## Argument Reference

* `resource_group_name` - Specifies the resource group where the Private DNS Zone (parent resource) exists.

* `zone_name` - Specifies the Private DNS Zone where the resource exists.

## Attributes Reference

* `id` - The Private DNS SOA Record ID.

* `name` - The name of the Private DNS SOA Record.

* `fqdn` - The FQDN of the Private DNS SOA Record.

* `ttl` - The Time To Live (TTL) of the Private DNS record in seconds.

* `email` - The email contact for the SOA record.

* `host_name` - The domain name of the authoritative name server for the SOA record.

* `expire_time` - The expire time for the SOA record.

* `minimum_ttl` - The minimum Time To Live for the SOA record. By convention, it is used to determine the negative caching duration.

* `refresh_time` - The refresh time for the SOA record.

* `retry_time` - The retry time for the SOA record.

* `serial_number` - The serial number for the SOA record.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS SOA Record.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network`: 2024-06-01
