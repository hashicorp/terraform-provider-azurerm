---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_policy_arc_machine_configuration_assignment"
description: |-
  Gets information about an existing Policy.
---

# Data Source: azurerm_policy_arc_machine_configuration_assignment

Use this data source to access information about an existing Guest Configuration Policy assignment on an Arc machine.

## Example Usage

```hcl
data "azurerm_policy_arc_machine_configuration_assignment" "example" {
  name                = "existing"
  resource_group_name = "existing"
  machine_name        = "existing"
}

output "id" {
  value = data.azurerm_policy_arc_machine_configuration_assignment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `machine_name` - (Required) The name of the HybridCompute machine which the Guest Configuration Assignment is applied to.

* `name` - (Required) The name of the Guest Configuration Assignment.

* `resource_group_name` - (Required) The name of the Resource Group where the Guest Configuration Assignment exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Guest Configuration Assignment.

* `assignment_hash` - Combined hash of the configuration package and parameters..

* `compliance_status` - A value indicating compliance status of the machine for the assigned Guest Configuration. Possible return values are `Compliant`, `NonCompliant` and `Pending`.

* `content_hash` - The sha256 content hash for the Guest Configuration package.

* `content_uri` - The content URI where the Guest Configuration package is stored.

* `last_compliance_status_checked` - Date and time, in RFC3339 format, when the machines compliance status was last checked.

* `latest_report_id` - The ID of the latest report for the Guest Configuration Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Guest Configuration Assignment.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.HybridCompute` - 2024-04-05
