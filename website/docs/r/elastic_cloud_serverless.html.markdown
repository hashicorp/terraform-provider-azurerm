---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_serverless"
description: |-
  Manages an Elastic Cloud Serverless project.
---

# azurerm_elastic_cloud_serverless

Manages an Elastic Cloud Serverless project using a `Microsoft.Elastic/monitors` resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "East US"
}

resource "azurerm_elastic_cloud_serverless" "example" {
  name                = "example-elastic-cloud-serverless"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  kind                        = "elastic-serverless-search"
  sku_name                    = "ess-consumption-2024_Monthly"
  project_type                = "Elasticsearch"
  configuration_type          = "GeneralPurpose"
  offer_id                    = "ec-azure-pp"
  term_id                     = "n7ja87drquhy"
  elastic_cloud_email_address = "alice@example.com"

  tags = {
    Environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Elastic Cloud Serverless project. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Elastic Cloud Serverless project should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Elastic Cloud Serverless project should exist. Changing this forces a new resource to be created.

* `configuration_type` - (Required) Specifies the Elastic project configuration type. Possible values are `GeneralPurpose`, `Vector`, `TimeSeries`, and `NotApplicable`. Changing this forces a new resource to be created.

* `elastic_cloud_email_address` - (Required) Specifies the email address of the Elastic Cloud user associated with the authenticated Azure principal. Changing this forces a new resource to be created.

* `kind` - (Required) Specifies the Serverless resource kind associated with the selected Elastic project and Marketplace plan. Changing this forces a new resource to be created.

* `offer_id` - (Required) Specifies the Elastic Marketplace offer ID associated with the selected SKU and project type. Changing this forces a new resource to be created.

* `project_type` - (Required) Specifies the Elastic project type. Possible values are `Elasticsearch`, `Observability`, and `Security`. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the name of the Elastic Marketplace SKU. Changing this forces a new resource to be created.

* `term_id` - (Required) Specifies the Elastic Marketplace term ID associated with the selected SKU. Changing this forces a new resource to be created.

* `generate_api_key` - (Optional) Specifies whether Elastic should generate an API key for the project. Defaults to `false`. Changing this forces a new resource to be created.

* `monitoring_enabled` - (Optional) Specifies whether monitoring is enabled for the Elastic Cloud Serverless project. Defaults to `true`. Changing this forces a new resource to be created.

* `plan_id` - (Optional) Specifies the Elastic Marketplace plan ID. Defaults to `ess-consumption-2024`. Changing this forces a new resource to be created.

* `publisher_id` - (Optional) Specifies the Elastic Marketplace publisher ID. Defaults to `elastic`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elastic Cloud Serverless project.

~> **Note:** `kind`, `sku_name`, `project_type`, `configuration_type`, `offer_id`, `plan_id`, `publisher_id`, and `term_id` describe one Elastic Marketplace plan and must be supplied as a compatible combination.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic Cloud Serverless project.

* `elastic_cloud_deployment_id` - The ID of the deployment within Elastic Cloud.

* `elasticsearch_service_url` - The URL of the Elasticsearch service associated with this project.

* `kibana_service_url` - The URL of the Kibana service associated with this project.

* `kibana_sso_uri` - The URI used for SSO to the Kibana service associated with this project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Elastic Cloud Serverless project.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Cloud Serverless project.
* `update` - (Defaults to 1 hour) Used when updating the Elastic Cloud Serverless project.
* `delete` - (Defaults to 1 hour) Used when deleting the Elastic Cloud Serverless project.

## Import

An Elastic Cloud Serverless project can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_cloud_serverless.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Elastic/monitors/monitor1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Elastic` - 2025-06-01
