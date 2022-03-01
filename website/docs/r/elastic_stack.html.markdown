---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_stack"
description: |-
  Manages an Elastic Stack in Elastic Cloud.
---

# azurerm_elastic_stack

Manages an Elastic Stack in Elastic Cloud.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_elastic_stack" "test" {
  name                        = "example-elastic-stack"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-monthly-consumption_Monthly"
  elastic_cloud_email_address = "user@example.com"
}
```

## Arguments Reference

The following arguments are supported:

* `elastic_cloud_email_address` - (Required) Specifies the Email Address which should be associated with this Elastic Stack account. Changing this forces a new Elastic Stack to be created.

* `location` - (Required) The Azure Region where the Elastic Stack should exist. Changing this forces a new Elastic Stack to be created.

* `name` - (Required) The name which should be used for this Elastic Stack. Changing this forces a new Elastic Stack to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Elastic Stack should exist. Changing this forces a new Elastic Stack to be created.

* `sku_name` - (Required) Specifies the name of the SKU for this Elastic Stack. Changing this forces a new Elastic Stack to be created.

---

* `logs` - (Optional) A `logs` block as defined below.

* `monitoring_enabled` - (Optional) Specifies if the Elastic Stack should have monitoring configured? Defaults to `true`. Changing this forces a new Elastic Stack to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elastic Stack.

---

The `filtering_tag` block supports the following:

* `action` - (Required) Specifies the type of action which should be taken when the Tag matches the `name` and `value`. Possible values are `Exclude` and `Include`.

* `name` - (Required) Specifies the name (key) of the Tag which should be filtered.

* `value` - (Required) Specifies the value of the Tag which should be filtered.

---

The `logs` block supports the following:

* `filtering_tag` - (Optional) A list of `filtering_tag` blocks as defined above.

* `send_activity_logs` - (Optional) Specifies if the Azure Activity Logs should be sent to the Elastic Stack cluster. Defaults to `false`.

* `send_azuread_logs` - (Optional) Specifies if the AzureAD Logs should be sent to the Elastic Stack cluster. Defaults to `false`.

* `send_subscription_logs` - (Optional) Specifies if the Azure Subscription Logs should be sent to the Elastic Stack cluster. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Elastic Stack.

* `elastic_cloud_deployment_id` - The ID of the Deployment within Elastic Cloud.

* `elastic_cloud_sso_default_url` - The Default URL used for Single Sign On (SSO) to Elastic Cloud.

* `elastic_cloud_user_id` - The ID of the User Account within Elastic Cloud.

* `elasticsearch_service_url` - The URL to the Elasticsearch Service associated with this Elastic Stack.

* `kibana_service_url` - The URL to the Kibana Dashboard associated with this Elastic Stack.

* `kibana_sso_uri` - The URI used for SSO to the Kibana Dashboard associated with this Elastic Stack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Elastic Stack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Elastic Stack.
* `update` - (Defaults to 30 minutes) Used when updating the Elastic Stack.
* `delete` - (Defaults to 30 minutes) Used when deleting the Elastic Stack.

## Import

Elastic Stack's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_stack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Elastic/monitors/monitor1
```
