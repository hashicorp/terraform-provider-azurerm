---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_aws_cloud_trail"
description: |-
  Manages a AWS CloudTrail Data Connector.
---

# azurerm_sentinel_data_connector_aws_cloud_trail

Manages a AWS CloudTrail Data Connector.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_data_connector_aws_cloud_trail" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  aws_role_arn               = "arn:aws:iam::000000000000:role/role1"
}
```

## Arguments Reference

The following arguments are supported:

* `aws_role_arn` - (Required) The ARN of the AWS CloudTrail role, which is connected to this AWS CloudTrail Data Connector.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this AWS CloudTrail Data Connector resides in. Changing this forces a new AWS CloudTrail Data Connector to be created.

* `name` - (Required) The name which should be used for this AWS CloudTrail Data Connector. Changing this forces a new AWS CloudTrail Data Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the AWS CloudTrail Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the AWS CloudTrail Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the AWS CloudTrail Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the AWS CloudTrail Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the AWS CloudTrail Data Connector.

## Import

AWS CloudTrail Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_aws_cloud_trail.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
