---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_email_communication_service_domain"
description: |-
  Manages an Email Communication Service Domain.
---

# azurerm_email_communication_service_domain

Manages an Email Communication Service Domain.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_email_communication_service" "example" {
  name                = "example-emailcommunicationservice"
  resource_group_name = azurerm_resource_group.example.name
  data_location       = "United States"
}

resource "azurerm_email_communication_service_domain" "example" {
  name             = "AzureManagedDomain"
  email_service_id = azurerm_email_communication_service.example.id

  domain_management = "AzureManaged"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Email Communication Service resource. If `domain_management` is `AzureManaged`, the name must be `AzureManagedDomain`. Changing this forces a new Email Communication Service to be created.

* `email_service_id` - (Required) The resource ID of the Email Communication Service where the Domain belongs to. Changing this forces a new Email Communication Service to be created.

* `domain_management` - (Required) Describes how a Domains resource is being managed. Possible values are `AzureManaged`, `CustomerManaged`, `CustomerManagedInExchangeOnline`. Changing this forces a new Email Communication Service to be created.

---

* `user_engagement_tracking_enabled` - (Optional) Describes user engagement tracking is enabled or disabled. Defaults to `false`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Email Communication Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Communication Service.

* `from_sender_domain` - P2 sender domain that is displayed to the email recipients [RFC 5322].

* `mail_from_sender_domain` - P1 sender domain that is present on the email envelope [RFC 5321].

* `verification_records` - (Optional) An `verification_records` block as defined below.

An `verification_records` block supports the following arguments:

* `domain` - (Optional) An `domain` block as defined below.

* `dkim` - (Optional) An `dkim` block as defined below.

* `dkim2` - (Optional) An `dkim2` block as defined below.

* `dmarc` - (Optional) An `dmarc` block as defined below.

* `spf` - (Optional) An `spf` block as defined below.

An `domain` block supports the following arguments:

* `name` - Name of the DNS record.

* `ttl` - Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.

* `type` - Type of the DNS record. Example: TXT

* `value` - Value of the DNS record.

An `dkim` block supports the following arguments:

* `name` - Name of the DNS record.

* `ttl` - Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.

* `type` - Type of the DNS record. Example: TXT

* `value` - Value of the DNS record.

An `dkim2` block supports the following arguments:

* `name` - Name of the DNS record.

* `ttl` - Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.

* `type` - Type of the DNS record. Example: TXT

* `value` - Value of the DNS record.

An `dmarc` block supports the following arguments:

* `name` - Name of the DNS record.

* `ttl` - Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.

* `type` - Type of the DNS record. Example: TXT

* `value` - Value of the DNS record.

An `spf` block supports the following arguments:

* `name` - Name of the DNS record.

* `ttl` - Represents an expiry time in seconds to represent how long this entry can be cached by the resolver, default = 3600sec.

* `type` - Type of the DNS record. Example: TXT

* `value` - Value of the DNS record.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Communication Service Domain.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Communication Service Domain.
* `update` - (Defaults to 30 minutes) Used when updating the Email Communication Service Domain.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Communication Service Domain.

## Import

Communication Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_email_communication_service_domain.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/emailServices/emailCommunicationService1/domains/domain1
```
