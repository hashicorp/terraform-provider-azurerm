subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud"
sidebar_current: "docs-azurerm-resource-spring-cloud"
description: |-
  Manage Azure Spring Cloud Service instance.
---

# azurerm_spring_cloud

Manage Azure Spring Cloud Service instance. Azure Spring Cloud provides a managed service that enables Java developers to easily build and run Spring-boot based microservices on Azure with no code changes.
Within an azure spring cloud service, users can manage config server, manage multiple spring cloud apps and manage deployments

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Service resource. Changing this forces a new resource to be created.

* `resource_group` - (Required) The name of the resource group that contains the resource. You can obtain this value from the Azure Resource Manager API or the portal. Changing this forces a new resource to be created.

* `location` - (Optional) The GEO location of the resource. Changing this forces a new resource to be created. 

* `tags` - (Optional) Tags of the service which is a list of key value pairs that describe the resource. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `version` - Version of the Service

* `service_id` - ServiceInstanceEntity GUID which uniquely identifies a created resource

* `id` - Fully qualified resource Id for the resource.

* `name` - The name of the resource.
