---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_mx_record"
description: |-
  Gets information about an existing DNS MX Record.
---

# Data Source: azurerm_dns_mx_record

Use this data source to access information about an existing DNS MX Record within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

## Example Usage

```hcl
resource "azurerm_dns_mx_record" "example" {
  name                = "test"
  zone_name           = "test-zone"
  resource_group_name = "test-rg"
}

output "dns_mx_record_id" {
  value = data.azurerm_dns_mx_record.example.id
}
```

## Argument Reference

* `name` - The name of the DNS MX Record.

* `resource_group_name` - Specifies the resource group where the DNS Zone (parent resource) exists.

* `zone_name` - Specifies the DNS Zone where the resource exists.

## Attributes Reference

* `id` - The DNS MX Record ID.

* `fqdn` - The FQDN of the DNS MX Record.

* `ttl` - The Time To Live (TTL) of the DNS record in seconds.

* `record` - A list of values that make up the MX record. Each `record` block supports fields documented below.

* `tags` - A mapping of tags assigned to the resource.

---

The `record` block supports:

* `preference` - String representing the "preference” value of the MX records. Records with lower preference value take priority.

* `exchange` - The mail server responsible for the domain covered by the MX record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS MX Record.
