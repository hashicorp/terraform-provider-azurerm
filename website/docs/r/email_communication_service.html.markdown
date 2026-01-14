---
subcategory: "Communication"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_email_communication_service"
description: |-
  Manages an Email Communication Service.
---

# azurerm_email_communication_service

Manages an Email Communication Service.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Email Communication Service resource. Changing this forces a new Email Communication Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Email Communication Service should exist. Changing this forces a new Email Communication Service to be created.

---

* `data_location` - (Required) The location where the Email Communication service stores its data at rest. Possible values are `Africa`, `Asia Pacific`, `Australia`, `Brazil`, `Canada`, `Europe`, `France`, `Germany`, `India`, `Japan`, `Korea`, `Norway`, `Switzerland`, `UAE`, `UK` `usgov` and `United States`. Changing this forces a new Email Communication Service to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Email Communication Service.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Email Communication Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Email Communication Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Email Communication Service.
* `update` - (Defaults to 30 minutes) Used when updating the Email Communication Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Email Communication Service.

## Import

`azurerm_email_communication_service` resources can be imported using one of the following methods:

* The `terraform import` CLI command with an `id` string:

  ```shell
  terraform import azurerm_email_communication_service.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{emailServiceName}
  ```

* An `import` block with an `id` argument:
  
  ```hcl
  import {
    to = azurerm_email_communication_service.example
    id = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{emailServiceName}"
  }
  ```

* An `import` block with an `identity` argument:

  ```hcl
  import {
    to       = azurerm_email_communication_service.example
    identity = {
      TODO Resource Identity Format
    }
  }
  ```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Communication` - 2023-03-31
