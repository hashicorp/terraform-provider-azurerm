---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_monitor"
description: |-
  Manages a Elastic Monitor.
---

# azurerm_elastic_monitor

Manages a elastic Monitor.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "West US 2"
  location = "%s"
}
resource "azurerm_elastic_monitor" "test" {
  name                = "example_elastic_monitor"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  user_info {
    email_address = "abc@microsoft.com"
  }
  sku {
    name = "staging_Monthly"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Elastic Monitor. Changing this forces a new Elastic Monitor to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the elastic monitor should exist. Changing this forces a new elastic Monitor to be created.

* `location` - (Required) The Azure Region where the elastic Monitor should exist. Changing this forces a new elastic Monitor to be created.

* `sku` - (Required) A `sku` block as defined below.

* `user_info` - (Required) A `user_info` block as defined below.

---

* `monitoring_status` - (Optional) Flag specifying if the resource monitoring is enabled or disabled. Default value is `true`. Changing this forces a new elastic Monitor to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Elastic Monitor.

---

An `sku` block exports the following:

* `name` - (Required) Name of SKU of the monitor resource. Changing this forces a new elastic Monitor to be created.

---

An `user_info` block exports the following:

* `email_address` - (Required) Email of the user used to associate with elastic account. Changing this forces a new elastic Monitor to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the elastic Monitor.

* `elastic_properties` - The properties of Monitor. An `elastic_properties` is defined below.

---

An `elastic_properties` block exports the following:

* `elastic_cloud_user` - The properties of user associated with the elastic monitor. An `elastic_cloud_user` is defined below.

* `elastic_cloud_deployment` - The properties of elastic cloud deployment. An `elastic_cloud_deployment` is defined below.

---

An `elastic_cloud_user` block exports the following:

* `email_address` - Email of the Elastic User Account.

* `id` - User Id of the elastic account of the User.

* `elastic_cloud_sso_default_url` - Elastic cloud default dashboard sso URL of the Elastic user account.

---

An `elastic_cloud_deployment` block exports the following:

* `name` - Elastic deployment name.

* `deployment_id` - Elastic deployment Id.

* `azure_subscription_id` - Associated Azure subscription Id for the elastic deployment.

* `elasticsearch_region` - Region where Deployment at Elastic side took place.

* `elasticsearch_service_url` - Elasticsearch ingestion endpoint of the Elastic deployment.

* `kibana_service_url` - Kibana endpoint of the Elastic deployment.

* `kibana_sso_url` - Kibana dashboard sso URL of the Elastic deployment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the elastic Monitor.
* `read` - (Defaults to 5 minutes) Used when retrieving the elastic Monitor.
* `update` - (Defaults to 30 minutes) Used when updating the elastic Monitor.
* `delete` - (Defaults to 30 minutes) Used when deleting the elastic Monitor.

## Import

elastic Monitors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_elastic_monitor.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Elastic/monitors/monitor1
```
