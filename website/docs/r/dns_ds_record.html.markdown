---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_ds_record"
description: |-
  Manages a DNS DS Record.
---

# azurerm_dns_ds_record

Enables you to manage DNS DS Records within Azure DNS.

~> **Note:** [The Azure DNS API has a throttle limit of 500 read (GET) operations per 5 minutes](https://docs.microsoft.com/azure/azure-resource-manager/management/request-limits-and-throttling#network-throttling) - whilst the default read timeouts will work for most cases - in larger configurations you may need to set a larger [read timeout](https://www.terraform.io/language/resources/syntax#operation-timeouts) then the default 5min. Although, we'd generally recommend that you split the resources out into smaller Terraform configurations to avoid the problem entirely.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_dns_ds_record" "example" {
  name        = "test"
  dns_zone_id = azurerm_dns_zone.example.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the DNS DS Record. Changing this forces a new resource to be created.

- `dns_zone_id` - (Required) Specifies the DNS Zone ID where the resource exists. Changing this forces a new resource to be created.

- `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

- `record` - (Required) A list of values that make up the DS record. Each `record` block supports fields documented below.

- `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `record` block supports:

- `algorithm` - (Required) Identifies the algorithm used to produce a legitimate signature.

- `key_tag` - (Required) Contains the tag value of the DNSKEY Resource Record that validates this signature.

- `digest_type` - (Required) Identifies the algorithm used to construct the digest.

- `digest_value` - (Required) A cryptographic hash value of the referenced DNSKEY Record.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The DNS DS Record ID.

- `fqdn` - The FQDN of the DNS DS Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS DS Record.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS DS Record.

* `update` - (Defaults to 30 minutes) Used when updating the DNS DS Record.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS DS Record.

## Import

DS records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_ds_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/DS/myrecord1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2023-07-01-preview
