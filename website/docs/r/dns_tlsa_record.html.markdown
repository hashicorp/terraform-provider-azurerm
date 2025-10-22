---
subcategory: "DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dns_tlsa_record"
description: |-
  Manages a DNS TLSA Record.
---

# azurerm_dns_tlsa_record

Enables you to manage DNS TLSA Records within Azure DNS.

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

resource "azurerm_dns_tlsa_record" "example" {
  name        = "test"
  dns_zone_id = azurerm_dns_zone.example.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370c66fd4a0673ce1b62e76b819835dabb20702e4497cb10affe46e8135381e7"
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the DNS TLSA Record. Changing this forces a new resource to be created.

- `dns_zone_id` - (Required) Specifies the DNS Zone ID where the resource exists. Changing this forces a new resource to be created.

- `ttl` - (Required) The Time To Live (TTL) of the DNS record in seconds.

- `record` - (Required) A list of values that make up the TLSA record. Each `record` block supports fields documented below.

- `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `record` block supports:

- `matching_type` - (Required) Matching Type of the TLSA record.

- `selector` - (Required) Selector of the TLSA record.

- `usage` - (Required) Usage of the TLSA record.

- `certificate_association_data` - (Required) Certificate data to be matched.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The DNS TLSA Record ID.

- `fqdn` - The FQDN of the DNS TLSA Record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DNS TLSA Record.

* `read` - (Defaults to 5 minutes) Used when retrieving the DNS TLSA Record.

* `update` - (Defaults to 30 minutes) Used when updating the DNS TLSA Record.

* `delete` - (Defaults to 30 minutes) Used when deleting the DNS TLSA Record.

## Import

TLSA records can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dns_tlsa_record.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/dnsZones/zone1/TLSA/myrecord1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Network` - 2023-07-01-preview
