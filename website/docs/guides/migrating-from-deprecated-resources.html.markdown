---
layout: "azurerm"
page_title: "Azure Provider: Migrating from Deprecated Resources Guide"
description: |-
  This page documents how to migrate from deprecated resources in the Azure Provider to their replacements.
---

# Azure Provider: Migrating from Deprecated Resources Guide

This guide shows how to migrate from a resource which has been deprecated or renamed to its replacement.

It's possible to migrate between the resources by updating your Terraform Configuration, removing the old state, and the importing the new resource in config.

In this guide, we'll assume we're migrating from the `azurerm_app_service` resource to the new `azurerm_linux_web_app` resource, but this should also be applicable for resources that have only been renamed where you can simply change the resource type name in your config.

Assuming we have the following Terraform Configuration:

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_app_service_plan" "example" {
  name                = "Example App Service Plan"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  kind                = "Linux"
  reserved            = true

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "Example App Service"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  app_service_plan_id = azurerm_app_service_plan.main.id

  site_config {
    dotnet_framework_version = "v4.0"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2019"
  }
}
```

We can update the Terraform Configuration to use the new resource by updating the resources to the new `azurerm_service_plan` and `azurerm_linux_web_app` schema:

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_service_plan" "example" {
  name                = "Example App Service Plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "example" {
  name                = "Example App Service"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_service_plan.example.location
  service_plan_id     = azurerm_service_plan.example.id

  site_config {
    # ...
  }
}
```

As the Terraform Configuration has been updated - we now need to update the State. We can view the items Terraform is tracking in its statefile using the `terraform state list` command, for example:

```bash
$ terraform state list
azurerm_app_service.example
azurerm_app_service_plan.example
azurerm_resource_group.example
```

In order to migrate from the old resource to the new resource we need to first remove the old resource from the state - and subsequently use Terraform's [import functionality](https://www.terraform.io/docs/import/index.html) to migrate to the new resource.

To import a resource in Terraform we first require its Resource ID - we can obtain this from the command-line via:

```shell
$ echo azurerm_app_service_plan.example.id | terraform console
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1
$ echo azurerm_app_service.example.id | terraform console
/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1
```

Next we can remove the existing resource using `terraform state rm` - for example:

```shell
$ terraform state rm azurerm_app_service.example azurerm_app_service_plan.example
Removed azurerm_autoscale_setting.example
Successfully removed 2 resource instance(s).
```

Now that the old resource has been removed from Terraform's Statefile we can now Import it into the Statefile as the new resource by running:

```text
terraform import [resourcename].[identifier] [resourceid]
```

For example:

```shell
$ terraform import azurerm_service_plan.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1
azurerm_monitor_autoscale_setting.test: Importing from ID "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1"...
azurerm_monitor_autoscale_setting.test: Import prepared!
  Prepared azurerm_monitor_autoscale_setting for import
azurerm_monitor_autoscale_setting.test: Refreshing state... [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/serverfarms/instance1]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.

$ terraform import azurerm_linux_web_app.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1
azurerm_monitor_autoscale_setting.test: Importing from ID "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1"...
azurerm_monitor_autoscale_setting.test: Import prepared!
  Prepared azurerm_monitor_autoscale_setting for import
azurerm_monitor_autoscale_setting.test: Refreshing state... [id=/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1]

Import successful!

The resources that were imported are shown above. These resources are now in
your Terraform state and will henceforth be managed by Terraform.
```

Once this has been done, running `terraform plan` should show no changes:

```shell
$ terraform plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

At this point, you've switched over to using the new resource and should be able to continue using Terraform as normal.
