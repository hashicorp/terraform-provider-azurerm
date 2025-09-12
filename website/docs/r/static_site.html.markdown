---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_static_site"
description: |-
  Manages a Static Site.
---

# azurerm_static_site

Manages an App Service Static Site.

-> **Note:** The `azurerm_static_site` resource is deprecated in favour of `azurerm_static_web_app` and will be removed in a future major release.

-> **Note:** After the Static Site is provisioned, you'll need to associate your target repository, which contains your web app, to the Static Site, by following the [Azure Static Site document](https://docs.microsoft.com/azure/static-web-apps/github-actions-workflow).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_static_site" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Static Web App. Changing this forces a new Static Web App to be created.

* `location` - (Required) The Azure Region where the Static Web App should exist. Changing this forces a new Static Web App to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Static Web App should exist. Changing this forces a new Static Web App to be created.

* `sku_tier` - (Optional) Specifies the SKU tier of the Static Web App. Possible values are `Free` or `Standard`. Defaults to `Free`.

* `sku_size` - (Optional) Specifies the SKU size of the Static Web App. Possible values are `Free` or `Standard`. Defaults to `Free`.

* `identity` - (Optional) An `identity` block as defined below.

* `app_settings` - (Optional) A key-value pair of App Settings.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) The Type of Managed Identity assigned to this Static Site resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A list of Managed Identity IDs which should be assigned to this Static Site resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Static Web App.

* `api_key` - The API key of this Static Web App, which is used for later interacting with this Static Web App from other clients, e.g. GitHub Action.
  
* `default_host_name` - The default host name of the Static Web App.

* `identity` - (Optional) An `identity` block as defined below which contains the Managed Service Identity information for this resource.

---

An `identity` block exports the following:

* `type` - (Required) The Type of Managed Identity assigned to this resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.

* `principal_id` - (Optional) The Principal ID associated with this Managed Service Identity.

-> **Note:** You can access the Principal ID via `azurerm_static_site.example.identity[0].principal_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Static Web App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Static Web App.
* `update` - (Defaults to 30 minutes) Used when updating the Static Web App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Static Web App.

## Import

Static Web Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_static_site.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Web/staticSites/my-static-site1
```
