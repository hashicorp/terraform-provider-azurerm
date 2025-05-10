---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_management_group_policy_remediation"
description: |-
  Manages an Azure Management Group Policy Remediation.
---

# azurerm_management_group_policy_remediation

Manages an Azure Management Group Policy Remediation.

## Example Usage

```hcl
resource "azurerm_management_group" "example" {
  display_name = "Example Management Group"
}

data "azurerm_policy_definition" "example" {
  display_name = "Allowed locations"
}

resource "azurerm_management_group_policy_assignment" "example" {
  name                 = "exampleAssignment"
  management_group_id  = azurerm_management_group.example.id
  policy_definition_id = data.azurerm_policy_definition.example.id
  parameters = jsonencode({
    "listOfAllowedLocations" = {
      "value" = ["East US"]
    }
  })
}

resource "azurerm_management_group_policy_remediation" "example" {
  name                 = "example"
  management_group_id  = azurerm_management_group.example.id
  policy_assignment_id = azurerm_management_group_policy_assignment.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Policy Remediation. Changing this forces a new resource to be created.

* `management_group_id` - (Required) The Management Group ID at which the Policy Remediation should be applied. Changing this forces a new resource to be created.

* `policy_assignment_id` - (Required) The ID of the Policy Assignment that should be remediated.

* `policy_definition_reference_id` - (Optional) The unique ID for the policy definition reference within the policy set definition that should be remediated. Required when the policy assignment being remediated assigns a policy set definition.

* `location_filters` - (Optional) A list of the resource locations that will be remediated.

* `failure_percentage` - (Optional) A number between 0.0 to 1.0 representing the percentage failure threshold. The remediation will fail if the percentage of failed remediation operations (i.e. failed deployments) exceeds this threshold.

* `parallel_deployments` - (Optional) Determines how many resources to remediate at any given time. Can be used to increase or reduce the pace of the remediation. If not provided, the default parallel deployments value is used.

* `resource_count` - (Optional) Determines the max number of resources that can be remediated by the remediation job. If not provided, the default resource count is used.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Policy Remediation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Remediation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Remediation.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Remediation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Remediation.

## Import

Policy Remediations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_management_group_policy_remediation.example /providers/Microsoft.Management/managementGroups/my-mgmt-group-id/providers/Microsoft.PolicyInsights/remediations/remediation1
```
