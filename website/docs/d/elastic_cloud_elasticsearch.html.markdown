---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_elasticsearch"
description: |- 
    Gets information about an existing Elasticsearch resource.

---

# Data Source: azurerm_elastic_cloud_elasticsearch

Use this data source to access information about an existing Elasticsearch resource.

## Example Usage

```hcl
data "azurerm_elastic_cloud_elasticsearch" "example" {
  name                = "my-elastic-search"
  resource_group_name = "example-resources"
}

output "elasticsearch_endpoint" {
  value = data.azurerm_elastic_cloud_elasticsearch.example.elasticsearch_service_url
}

output "kibana_endpoint" {
  value = data.azurerm_elastic_cloud_elasticsearch.example.kibana_service_url
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Elasticsearch resource.

* `resource_group_name` - (Required) The name of the resource group in which the Elasticsearch exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Elasticsearch.

* `elastic_cloud_deployment_id` - The ID of the Deployment within Elastic Cloud.

* `elastic_cloud_email_address` - The Email Address which is associated with this Elasticsearch account.

* `elastic_cloud_sso_default_url` - The Default URL used for Single Sign On (SSO) to Elastic Cloud.

* `elastic_cloud_user_id` - The ID of the User Account within Elastic Cloud.

* `elasticsearch_service_url` - The URL to the Elasticsearch Service associated with this Elasticsearch.

* `kibana_service_url` - The URL to the Kibana Dashboard associated with this Elasticsearch.

* `kibana_sso_uri` - The URI used for SSO to the Kibana Dashboard associated with this Elasticsearch.

* `location` - The Azure Region in which this Elasticsearch exists.

* `logs` - A `logs` block as defined below.

* `monitoring_enabled` - Specifies if monitoring is enabled on this Elasticsearch or not.

* `sku_name` - The name of the SKU used for this Elasticsearch.

* `tags` - A mapping of tags assigned to the Elasticsearch.

---

The `filtering_tag` block exports the following:

* `action` - The type of action which is taken when the Tag matches the `name` and `value`.

* `name` - The name (key) of the Tag which should be filtered.

* `value` - The value of the Tag which should be filtered.

---

The `logs` block exports the following:

* `filtering_tag` - A list of `filtering_tag` blocks as defined above.

* `send_activity_logs` - Should the Azure Activity Logs should be sent to the Elasticsearch cluster?

* `send_azuread_logs` - Should the AzureAD Logs should be sent to the Elasticsearch cluster?

* `send_subscription_logs` - Should the Azure Subscription Logs should be sent to the Elasticsearch cluster?

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Elasticsearch.
