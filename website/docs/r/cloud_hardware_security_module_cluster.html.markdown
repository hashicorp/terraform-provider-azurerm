---
subcategory: "Hardware Security Module"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cloud_hardware_security_module_cluster"
description: |-
  Manages a Cloud Hardware Security Module Cluster.
---

# azurerm_cloud_hardware_security_module_cluster

Manages a Cloud Hardware Security Module Cluster.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cloud_hardware_security_module_cluster" "example" {
  name                = "example-cloudhsm-cluster"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  tags = {
    environment = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Hardware Security Module Cluster. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Cloud Hardware Security Module Cluster should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Cloud Hardware Security Module Cluster should exist. Changing this forces a new resource to be created.

---

* `identity` - (Optional) An `identity` block as defined below.

* `domain_name_reuse` - (Optional) Specifies the scope for auto-generated domain name labels. Possible values are `TenantReuse`, `SubscriptionReuse`, `ResourceGroupReuse`, and `NoReuse`. Defaults to `TenantReuse`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Cloud Hardware Security Module Cluster.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cloud Hardware Security Module Cluster. The only supported value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cloud Hardware Security Module Cluster.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cloud Hardware Security Module Cluster.

* `hardware_security_module` - A list of `hardware_security_module` blocks as defined below.

---

A `hardware_security_module` block exports the following:

* `fqdn` - The fully qualified domain name of the HSM instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Cloud Hardware Security Module Cluster.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cloud Hardware Security Module Cluster.
* `update` - (Defaults to 1 hour) Used when updating the Cloud Hardware Security Module Cluster.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cloud Hardware Security Module Cluster.

## Import

Cloud Hardware Security Module Clusters can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cloud_hardware_security_module_cluster.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HardwareSecurityModules/cloudHsmClusters/cluster1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.HardwareSecurityModules` - 2025-03-31
