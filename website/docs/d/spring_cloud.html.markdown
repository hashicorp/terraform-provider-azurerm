subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud"
sidebar_current: "docs-azurerm-datasource-spring-cloud"
description: |-
  Gets information about an existing Spring Cloud Service
---

# Data Source: azurerm_spring_cloud

Use this data source to access information about an existing Spring Cloud Service.


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Spring Cloud Service resource.

* `resource_group_name` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal.


## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.

* `name` - The name of the resource.

* `location` - The GEO location of the resource.

* `resource_group` - Resource group of the resource.

* `service_id` - ServiceInstanceEntity GUID which uniquely identifies a created resource

* `version` - Version of the Spring Cloud Service

* `tags` - Tags of the service which is a list of key value pairs that describe the resource.
