---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_data_connector_aws_s3"
description: |-
  Manages a AWS S3 Data Connector.
---

# azurerm_sentinel_data_connector_aws_s3

Manages a AWS S3 Data Connector.

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

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

resource "azurerm_sentinel_data_connector_aws_s3" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  aws_role_arn               = "arn:aws:iam::000000000000:role/role1"
  destination_table          = "AWSGuardDuty"
  sqs_urls                   = ["https://sqs.us-east-1.amazonaws.com/000000000000/example"]
  depends_on                 = [azurerm_log_analytics_solution.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this AWS S3 Data Connector. Changing this forces a new AWS S3 Data Connector to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace that this AWS S3 Data Connector resides in. Changing this forces a new AWS S3 Data Connector to be created.

* `aws_role_arn` - (Required) The ARN of the AWS role, which is connected to this AWS CloudTrail Data Connector. See the [Azure document](https://docs.microsoft.com/azure/sentinel/connect-aws?tabs=s3#create-an-aws-assumed-role-and-grant-access-to-the-aws-sentinel-account) for details.

* `destination_table` - (Required) The name of the Log Analytics table that will store the ingested data.

* `sqs_urls` - (Required) Specifies a list of AWS SQS urls for the AWS S3 Data Connector.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the AWS S3 Data Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the AWS S3 Data Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the AWS S3 Data Connector.
* `update` - (Defaults to 30 minutes) Used when updating the AWS S3 Data Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the AWS S3 Data Connector.

## Import

AWS S3 Data Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_data_connector_aws_s3.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/dataConnectors/dc1
```
