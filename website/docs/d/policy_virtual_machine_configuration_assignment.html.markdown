---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_virtual_machine_configuration_assignment"
description: |-
  Get information about a Guest Configuration Policy.
---

# Data Source: azurerm_policy_virtual_machine_configuration_assignment

Use this data source to access information about an existing Guest Configuration Policy.

## Example Usage

```hcl
data "azurerm_policy_virtual_machine_configuration_assignment" "example" {
  name                 = "AzureWindowsBaseline"
  resource_group_name  = "example-RG"
  virtual_machine_name = "example-vm"
}

output "compliance_status" {
  value = data.azurerm_policy_virtual_machine_configuration_assignment.example.compliance_status
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Guest Configuration Assignment.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group where the Guest Configuration Assignment exists.

* `virtual_machine_name` - (Required) Only retrieve Policy Set Definitions from this Management Group.

## Attributes Reference

* `id` - The ID of the Guest Configuration Assignment.

* `content_hash` - The content hash for the Guest Configuration package.

* `content_uri` - The content URI where the Guest Configuration package is stored.

* `assignment_hash` - Combined hash of the configuration package and parameters.

* `compliance_status` - A value indicating compliance status of the machine for the assigned guest configuration. Possible return values are `Compliant`, `NonCompliant` and `Pending`.

* `last_compliance_status_checked` - Date and time, in RFC3339 format, when the machines compliance status was last checked.

* `latest_report_id` - The ID of the latest report for the guest configuration assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Guest Configuration Assignment.
