---
layout: "azurerm"
page_title: "Azure Provider: Migrating to a renamed resource"
description: |-
    This page documents how to migrate between two resources in the Azure Provider which have been renamed.

---

# Azure Provider: Migrating to a renamed resource

This guide shows how to migrate from a resource which has been deprecated to its replacement. The complete list of resources that have been deprecated in favour of others can be found below.

~> **Note:** The following resources have been Deprecated and will be removed in version 2.0 of the Azure Provider

!> **Note:** You must move off the deprecated resources to upgrade to version 2.0 of the Azure Provider.

| Old Name                                               | New Name                                           |
| ------------------------------------------------------ | -------------------------------------------------- |
| azurerm_autoscale_setting                              | azurerm_monitor_autoscale_setting                  |
| azurerm_connection_monitor                             | azurerm_network_connection_monitor                 |
| azurerm_ddos_protection_plan                           | azurerm_network_ddos_protection_plan               |
| azurerm_log_analytics_workspace_linked_service         | azurerm_log_analytics_linked_service               |
| azurerm_iot_dps_certificate                            | azurerm_iothub_dps_certificate                     |
| azurerm_iot_dps                                        | azurerm_iothub_dps                                 |
| azurerm_private_link_endpoint                          | azurerm_private_endpoint                           |
| azurerm_private_link_endpoint_connection (Data Source) | azurerm_private_endpoint_connection                |
| azurerm_recovery_network_mapping                       | azurerm_site_recovery_network_mapping              |
| azurerm_recovery_replicated_vm                         | azurerm_site_recovery_replicated_vm                |
| azurerm_recovery_services_fabric                       | azurerm_site_recovery_fabric                       |
| azurerm_recovery_services_protected_vm                 | azurerm_backup_protected_vm                        |
| azurerm_recovery_services_protection_container         | azurerm_site_recovery_protection_container         |
| azurerm_recovery_services_protection_container_mapping | azurerm_site_recovery_protection_container_mapping |
| azurerm_recovery_services_protection_policy_vm         | azurerm_backup_policy_vm                           |
| azurerm_recovery_services_replication_policy           | azurerm_site_recovery_replication_policy           |

## Migrating to a renamed resource

As the Schema's for each resource are the same at this time - it's possible to migrate between the resources by updating your Terraform Configuration and updating the statefile.

In this guide, we'll assume we're migrating from the `azurerm_autoscale_setting` resource to the new `azurerm_monitor_autoscale_setting` resource, but this should be applicable for any of the resources listed below.

Assuming we have the following Terraform Configuration:

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_machine_scale_set" "example" {
  # ...
}

resource "azurerm_autoscale_setting" "example" {
  name                = "myAutoscaleSetting"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  target_resource_id  = azurerm_virtual_machine_scale_set.example.id
  # ...
}
```

We can update the Terraform Configuration to use the new resource by updating the name from `azurerm_autoscale_setting` to `azurerm_monitor_autoscale_setting`:

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_machine_scale_set" "example" {
  # ...
}

resource "azurerm_monitor_autoscale_setting" "example" {
  name                = "myAutoscaleSetting"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  target_resource_id  = azurerm_virtual_machine_scale_set.example.id
  # ...
}
```

As the Terraform Configuration has been updated - we now need to update the State. We can view the items Terraform is tracking in its Statefile using the `terraform state list` command, for example:

```bash
$ terraform state list
azurerm_autoscale_setting.example
azurerm_resource_group.example
azurerm_virtual_machine.example
```

In order to migrate from the old resource to the new resource we need to first remove the old resource from the state - and subsequently use Terraform's [import functionality](https://www.terraform.io/docs/import/index.html) to migrate to the new resource.

To import a resource in Terraform we first require its Resource ID - we can obtain this from the command-line via:

```shell
$ echo azurerm_autoscale_setting.example.id | terraform console
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/autoscalesettings/setting1
```

Next we can remove the existing resource using `terraform state rm` - for example:

```shell
$ terraform state rm azurerm_autoscale_setting.example
Removed azurerm_autoscale_setting.example
Successfully removed 1 resource instance(s).
```

Now that the old resource has been removed from Terraform's Statefile we can now Import it into the Statefile as the new resource by running:

```
$ terraform import [resourcename].[identifier] [resourceid]
```

For example:

```shell
$ terraform import azurerm_monitor_autoscale_setting.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/autoscalesettings/setting1
azurerm_monitor_autoscale_setting.test: Importing from ID "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/autoscalesettings/setting1"...
azurerm_monitor_autoscale_setting.test: Import prepared!
  Prepared azurerm_monitor_autoscale_setting for import
azurerm_monitor_autoscale_setting.test: Refreshing state... [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/microsoft.insights/autoscalesettings/setting1]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Once this has been done, running `terraform plan` should show no changes:

```shell
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

At this point, you've switched over to using the newly renamed resource and should be able to continue using Terraform as normal.
