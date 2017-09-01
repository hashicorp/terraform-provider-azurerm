---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service"
sidebar_current: "docs-azurerm-resource-app-service"
description: |-
  Create an App Service component.
---

# azurerm\_app\_service

Create an App Service component.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_app_service" "test" {
    name                = "api-appservice-pro"
    location            = "West Europe"
    resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service Plan component. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service Plan component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Optional) The resource ID of the app service plan.

* `always_on` - (Optional) Alows the app to be loaded all the time.

* `skip_dns_registration` - (Optional) If true, DNS registration is skipped.

* `skip_custom_domain_verification` - (Optional) If true, custom (non .azurewebsites.net) domains associated with web app are not verified.

* `force_dns_registration` - (Optional) If true, web app hostname is force registered with DNS.

* `ttl_in_seconds` - (Optional) Time to live in seconds for web app's default domain name.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service component.

## Import

App Service instances can be imported using the `resource id`, e.g.

```
terraform import azurerm_app_service.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1
```
