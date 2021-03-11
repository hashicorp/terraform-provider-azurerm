---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_policy"
description: |-
  Manages a Policy within a Dev Test Policy Set.
---

# azurerm_dev_test_policy

Manages a Policy within a Dev Test Policy Set.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_test_lab" "example" {
  name                = "example-devtestlab"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    "Sydney" = "Australia"
  }
}

resource "azurerm_dev_test_policy" "example" {
  name                = "LabVmCount"
  policy_set_name     = "default"
  lab_name            = azurerm_dev_test_lab.example.name
  resource_group_name = azurerm_resource_group.example.name
  fact_data           = ""
  threshold           = "999"
  evaluator_type      = "MaxValuePolicy"

  tags = {
    "Acceptance" = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Dev Test Policy. Possible values are `GalleryImage`, `LabPremiumVmCount`, `LabTargetCost`, `LabVmCount`, `LabVmSize`, `UserOwnedLabPremiumVmCount`, `UserOwnedLabVmCount` and `UserOwnedLabVmCountInSubnet`. Changing this forces a new resource to be created.

* `policy_set_name` - (Required) Specifies the name of the Policy Set within the Dev Test Lab where this policy should be created. Changing this forces a new resource to be created.

* `lab_name` - (Required) Specifies the name of the Dev Test Lab in which the Policy should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Dev Test Lab resource exists. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the Dev Test Lab exists. Changing this forces a new resource to be created.

* `description` - (Optional) A description for the Policy.

* `evaluator_type` - (Required) The Evaluation Type used for this Policy. Possible values include: 'AllowedValuesPolicy', 'MaxValuePolicy'. Changing this forces a new resource to be created.

* `threshold` - (Required) The Threshold for this Policy.

* `fact_data` - (Optional) The Fact Data for this Policy.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Dev Test Policy.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the DevTest Policy.
* `update` - (Defaults to 30 minutes) Used when updating the DevTest Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the DevTest Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the DevTest Policy.

## Import

Dev Test Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_policy.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/policysets/default/policies/policy1
```
