---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
description: |-
  Gets information about an existing Automation Account.
---

# Data Source: azurerm_automation_account

Use this data source to access information about an existing Automation Account.

## Example Usage

```hcl
data "azurerm_automation_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

output "automation_account_id" {
  value = data.azurerm_automation_account.example.id
}
```

## Argument Reference

* `name` - The name of the Automation Account.

* `resource_group_name` - Specifies the name of the Resource Group where the Automation Account exists.

## Attributes Reference

* `id` - The ID of the Automation Account

* `primary_key` - The Primary Access Key for the Automation Account.

* `secondary_key` - The Secondary Access Key for the Automation Account.

* `identity` - (Optional) An `identity` block as defined below.

* `endpoint` - The Endpoint for this Automation Account.

* `hybrid_service_url` - The URL of automation hybrid service which is used for hybrid worker on-boarding With this Automation Account.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Automation Account.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Automation Account.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Automation Account.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Automation Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Automation Account.
