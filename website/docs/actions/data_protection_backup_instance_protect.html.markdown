---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_instance_protect"
description: |-
  Sets the protection state of a Data Protection Backup Instance.
---

# Action: azurerm_data_protection_backup_instance_protect

~> **Note:** `azurerm_data_protection_backup_instance_protect` is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Changes the Protection state of a Data Protection Backup Instance to the specified value.

## Example Usage

```terraform
resource "azurerm_data_protection_backup_instance_postgresql_flexible_server" "example" {
  # ... Backup Instance configuration
}

resource "terraform_data" "example" {
  input = azurerm_data_protection_backup_instance_postgresql_flexible_server.example.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_data_protection_backup_instance_protect.stop_protection]
    }
  }
}



action "azurerm_data_protection_backup_instance_protect" "stop_protection" {
  config {
    backup_instance_id = azurerm_data_protection_backup_instance_postgresql_flexible_server.example.id
    protect_action     = "stop_protection"
  }
}

```

## Argument Reference

This action supports the following arguments:

* `backup_instance_id` - (Required) The ID of the data protection backup instance on which to perform the action.

* `protect_action` - (Required) The protect state action to take on this backup instance. Possible values include `stop_protection`,`resume_protection`, `suspend_backups`, and `resume_backups`.
