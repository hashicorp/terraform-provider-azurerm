---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_tlsa_record"
description: |-
  Gets information about an existing DNS TLSA Record.
---

# Data Source: azurerm_dns_tlsa_record

Use this data source to access information about an existing DNS TLSA Record within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

## Example Usage

```hcl
data "azurerm_dns_tlsa_record" "example" {
  name        = "test"
  dns_zone_id = azurerm_dns_zone.example.id
}

output "dns_tlsa_record_id" {
  value = data.azurerm_dns_tlsa_record.example.id
}
```

## Argument Reference

* `name` - The name of the DNS TLSA Record.

* `dns_zone_id` - Specifies the DNS Zone ID where the resource exists.

## Attributes Reference

* `id` - The DNS TLSA Record ID.

* `fqdn` - The FQDN of the DNS TLSA Record.

* `ttl` - The Time To Live (TTL) of the DNS record in seconds.

* `record` - A list of values that make up the TLSA record. Each `record` block supports fields documented below.

* `tags` - A mapping of tags assigned to the resource.

---

The `record` block supports:

* `matching_type` - Matching Type of the TLSA record.

* `selector` - Selector of the TLSA record.

* `usage` - Usage of the TLSA record.

* `certificate_association_data` - Certificate data to be matched.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS SRV Record.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Network` - 2023-07-01-preview
