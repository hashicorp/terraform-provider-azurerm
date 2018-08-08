---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_public_ip"
sidebar_current: "docs-azurerm-datasource-public-ip-x"
description: |-
  Retrieves information about the specified public IP address.

---

# Data Source: azurerm_public_ip

Use this data source to access the properties of an existing Azure Public IP Address.

## Example Usage

```hcl
data "azurerm_public_ip" "example" {
  name                = "example-publicip"
  resource_group_name = "example-resources"
}

output "public_ip_address" {
  value = "${data.azurerm_public_ip.example.ip_address}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the public IP address.
* `resource_group_name` - (Required) Specifies the name of the resource group.


## Attributes Reference

* `domain_name_label` - The label for the Domain Name.
* `idle_timeout_in_minutes` - Specifies the timeout for the TCP idle connection.
* `fqdn` - Fully qualified domain name of the A DNS record associated with the public IP. This is the concatenation of the domainNameLabel and the regionalized DNS zone.
* `ip_address` - The IP address value that was allocated.
* `tags` - A mapping of tags to assigned to the resource.
