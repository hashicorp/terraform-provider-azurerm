---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_email_communication_service_domain_sender_username"
description: |-
  Manages an Email Communication Service Domain Sender Username.
---

# azurerm_email_communication_service_domain_sender_username

Manages an Email Communication Service Domain Sender Username.

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
  name              = "AzureManagedDomain"
  email_service_id  = azurerm_email_communication_service.example.id
  domain_management = "AzureManaged"
}

resource "azurerm_email_communication_service_domain_sender_username" "example" {
  name                    = "example-su"
  email_service_domain_id = azurerm_email_communication_service_domain.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Email Communication Service Domain Sender Username resource. Changing this forces a new resource to be created.

* `email_service_domain_id` - (Required) The ID of the Email Communication Service Domain resource. Changing this forces a new resource to be created.

* `display_name` - (Optional) The display name for the Email Communication Service Domain Sender Username resource.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Communication Service Domain Sender Username.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Communication Service Domain Sender Username.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Communication Service Domain Sender Username.
* `update` - (Defaults to 30 minutes) Used when updating the Email Communication Service Domain Sender Username.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Communication Service Domain Sender Username.

## Import

Communication Service Domain Sender Usernames can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_email_communication_service_domain_sender_username.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Communication/emailServices/service1/domains/domain1/senderUsernames/username1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Communication` - 2023-03-31
