---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_cluster_arc_setting"
description: |-
  Gets information about an existing Azure Stack HCI Cluster Arc Setting.
---

# Data Source: azurerm_stack_hci_cluster_arc_setting

Use this data source to access information about an existing Azure Stack HCI Cluster Arc Setting instance.

## Example Usage

```hcl
data "azurerm_stack_hci_cluster" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

data "azurerm_stack_hci_cluster_arc_setting" "example" {
  name                 = "default"
  stack_hci_cluster_id = data.azurerm_stack_hci_cluster.example.id
}


output "id" {
  value = data.azurerm_stack_hci_cluster_arc_setting.example.id
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Stack HCI Cluster Arc Setting.

* `stack_hci_cluster_id` - (Required) The name of the Azure Stack HCI Cluster where the Arc Setting exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Stack HCI Cluster Arc Setting.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Cluster Arc Setting.
