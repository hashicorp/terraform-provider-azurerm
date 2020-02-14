---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_remediation"
description: |-
  Gets information about an existing Policy Remediation
---

# Data Source: azurerm_policy_insights_remediation

Use this data source to access information about an existing Policy Remediation.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the Policy Remediation.

* `scope` - (Required) Specifies the Scope of the Policy Remediation, which must be a Resource ID (such as Subscription e.g. `/subscriptions/00000000-0000-0000-0000-000000000000` or a Resource Group e.g.`/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup` or a specified Resource e.g. `/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Compute/virtualMachines/myVM`) or a Management Group (e.g. `/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000`).


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Policy Remediation.
