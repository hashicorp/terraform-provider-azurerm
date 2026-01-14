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

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Email Communication Service Domain Sender Username resource. Changing this forces a new resource to be created.

* `email_service_domain_id` - (Required) The ID of the Email Communication Service Domain resource. Changing this forces a new resource to be created.

* `display_name` - (Optional) The display name for the Email Communication Service Domain Sender Username resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Communication Service Domain Sender Username.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Communication Service Domain Sender Username.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Communication Service Domain Sender Username.
* `update` - (Defaults to 30 minutes) Used when updating the Email Communication Service Domain Sender Username.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Communication Service Domain Sender Username.

## Import

`azurerm_email_communication_service_domain_sender_username` resources can be imported using one of the following methods:

* The `terraform import` CLI command with an `id` string:

  ```shell
  terraform import azurerm_email_communication_service_domain_sender_username.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{emailServiceName}/domains/{domainName}/senderUsernames/{senderUsernameName}
  ```

* An `import` block with an `id` argument:
  
  ```hcl
  import {
    to = azurerm_email_communication_service_domain_sender_username.example
    id = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{emailServiceName}/domains/{domainName}/senderUsernames/{senderUsernameName}"
  }
  ```

* An `import` block with an `identity` argument:

  ```hcl
  import {
    to       = azurerm_email_communication_service_domain_sender_username.example
    identity = {
      TODO Resource Identity Format
    }
  }
  ```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Communication` - 2023-03-31
