---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_communication_service_email_domain_association"
description: |-
  Manages a communication service email domain association.
---

# azurerm_communication_service_email_domain_association

Manages a communication service email domain association.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "group1"
  location = "West Europe"
}

resource "azurerm_communication_service" "example" {
  name                = "CommunicationService1"
  resource_group_name = azurerm_resource_group.example.name
  data_location       = "United States"
}

resource "azurerm_email_communication_service" "example" {
  name                = "emailCommunicationService1"
  resource_group_name = azurerm_resource_group.example.name
  data_location       = "United States"
}

resource "azurerm_email_communication_service_domain" "example" {
  name             = "AzureManagedDomain"
  email_service_id = azurerm_email_communication_service.example.id

  domain_management = "AzureManaged"
}

resource "azurerm_communication_service_email_domain_association" "example" {
  communication_service_id = azurerm_communication_service.example.id
  email_service_domain_id  = azurerm_email_communication_service_domain.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `communication_service_id` - (Required) The ID of the Communication Service. Changing this forces a new communication service email domain association to be created.

* `email_service_domain_id` - (Required) The ID of the EMail Service Domain. Changing this forces a new communication service email domain association to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the communication service email domain association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the communication service email domain association.
* `read` - (Defaults to 5 minutes) Used when retrieving the communication service email domain association.
* `delete` - (Defaults to 5 minutes) Used when deleting the communication service email domain association.

## Import

Communication service email domain association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_communication_service_email_domain_association.example "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/communicationServices/communicationService1|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/emailServices/emailCommunicationService1/domains/domain1"
```
