---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_three_tier_virtual_instance"
description: |-
  Lists SAP Three Tier Virtual Instance resources.
---

# List resource: azurerm_workloads_sap_three_tier_virtual_instance

Lists SAP Three Tier Virtual Instance resources.

## Example Usage

### List all SAP Three Tier Virtual Instances in the subscription

```hcl
list "azurerm_workloads_sap_three_tier_virtual_instance" "example" {
  provider = azurerm
  config {}
}
```

### List all SAP Three Tier Virtual Instances in a specific resource group

```hcl
list "azurerm_workloads_sap_three_tier_virtual_instance" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported for each SAP Three Tier Virtual Instance:

* `id` - The ID of the SAP Virtual Instance.

* `name` - The name of the SAP Virtual Instance.

* `location` - The Azure Region where the SAP Virtual Instance exists.

* `resource_group_name` - The name of the Resource Group where the SAP Virtual Instance exists.

* `app_location` - The Geo-Location where the SAP system is to be deployed.

* `environment` - The environment type for the SAP Virtual Instance.

* `identity` - An `identity` block as defined below.

* `managed_resource_group_name` - The name of the managed Resource Group for the SAP Virtual Instance.

* `managed_resources_network_access_type` - The network access type for managed resources in the SAP Virtual Instance.

* `sap_fqdn` - The fully qualified domain name for the SAP system.

* `sap_product` - The SAP Product type.

* `three_tier_configuration` - A `three_tier_configuration` block as defined in the [azurerm_workloads_sap_three_tier_virtual_instance](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/workloads_sap_three_tier_virtual_instance) resource documentation.

* `tags` - A mapping of tags assigned to the SAP Virtual Instance.

---

An `identity` block exports the following:

* `type` - The Type of Managed Identity assigned to the SAP Virtual Instance.

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned to the SAP Virtual Instance.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Workloads` - 2024-09-01
