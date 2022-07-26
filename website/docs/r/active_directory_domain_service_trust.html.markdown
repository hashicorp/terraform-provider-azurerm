---
subcategory: "Active Directory Domain Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_active_directory_domain_service_trust"
description: |-
  Manages a Active Directory Domain Service Trust.
---

# azurerm_active_directory_domain_service_trust

Manages a Active Directory Domain Service Trust.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_active_directory_domain_service" "example" {
  name                = "example-ds"
  resource_group_name = "example-rg"
}

resource "azurerm_active_directory_domain_service_trust" "example" {
  name                   = "example-trust"
  domain_service_id      = data.azurerm_active_directory_domain_service.example.id
  trusted_domain_fqdn    = "example.com"
  trusted_domain_dns_ips = ["10.1.0.3", "10.1.0.4"]
  password               = "Password123"
}
```

## Arguments Reference

The following arguments are supported:

* `domain_service_id` - (Required) The ID of the Active Directory Domain Service. Changing this forces a new Active Directory Domain Service Trust to be created.

* `name` - (Required) The name which should be used for this Active Directory Domain Service Trust. Changing this forces a new Active Directory Domain Service Trust to be created.

* `password` - (Required) The password of the inbound trust set in the on-premise Active Directory Domain Service.

* `trusted_domain_dns_ips` - (Required) Specifies a list of DNS IPs that are used to resolve the on-premise Active Directory Domain Service.

* `trusted_domain_fqdn` - (Required) The FQDN of the on-premise Active Directory Domain Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Active Directory Domain Service Trust.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Active Directory Domain Service Trust.
* `read` - (Defaults to 5 minutes) Used when retrieving the Active Directory Domain Service Trust.
* `update` - (Defaults to 30 minutes) Used when updating the Active Directory Domain Service Trust.
* `delete` - (Defaults to 30 minutes) Used when deleting the Active Directory Domain Service Trust.

## Import

Active Directory Domain Service Trusts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_active_directory_domain_service_trust.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AAD/domainServices/DomainService1/trusts/trust1
```
