---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_extension"
description: |-
  Manages an Azure Stack HCI Extension.
---

# azurerm_stack_hci_extension

Manages an Azure Stack HCI Extension.

## Example Usage

```hcl
data "azurerm_stack_hci_cluster" "example" {
  name                = "aro-cluster"
  resource_group_name = "example-resources"
}

data "azurerm_stack_hci_cluster_arc_setting" "example" {
  name                 = "default"
  stack_hci_cluster_id = data.azurerm_stack_hci_cluster.example.id
}

resource "azurerm_stack_hci_extension" "example" {
  name                       = "example-shce"
  arc_setting_id             = data.azurerm_stack_hci_cluster_arc_setting.example.id
  publisher                  = "Microsoft.EnterpriseCloud.Monitoring"
  type                       = "MicrosoftMonitoringAgent"
  auto_upgrade_minor_version = true
  automatic_upgrade_enabled  = false
  force_update_tag           = "1"
  type_handler_version       = "1.22.0"

  protected_setting = <<PROTECTED_SETTING
{
	"workspaceKey": "xxxxx"
}
PROTECTED_SETTING

  setting = <<SETTING
{
	"workspaceId": "00000000-0000-0000-0000-000000000000"
}
SETTING
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Extension. Changing this forces a new resource to be created.

* `arc_setting_id` - (Required) The ID of the Azure Stack HCI Cluster Arc Setting. Changing this forces a new resource to be created.

* `publisher` - (Required) The name of the extension handler publisher, such as `Microsoft.Azure.Monitor`. Changing this forces a new resource to be created.

* `type` - (Required) Specifies the type of the extension. For example `CustomScriptExtension` or `AzureMonitorLinuxAgent`. Changing this forces a new resource to be created.

-> **NOTE:** The `Publisher` and `Type` of Azure Stack HCI Extensions can be found in this [Azure document](https://learn.microsoft.com/en-us/azure/azure-arc/servers/manage-vm-extensions#windows-extensions).

* `automatic_upgrade_enabled` - (Optional) Indicates whether the extension should be automatically upgraded by the platform if there is a newer version available. Possible values are `true` and `false`.

* `protected_settings` - (Optional) Json formatted protected settings for the extension.

* `settings` - (Optional) Json formatted public settings for the extension.

* `type_handler_version` - (Optional) Specifies the version of the script handler.

-> **NOTE:** `type_handler_version` cannot be set when `automatic_upgrade_enabled` is set to `true`.

-> **NOTE:** Possible values for `type_handler_version` can be found using the Azure CLI, via:

```shell
az vm extension image list --publisher Microsoft.Azure.Monitor -n AzureMonitorWindowsAgent --location westus -o table
```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Stack HCI Extension.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Extension.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Extension.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Extension.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Extension.

## Import

Azure Stack HCI Extension can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_extension.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/clusters/cluster1/arcSettings/default/extensions/extension1
```
