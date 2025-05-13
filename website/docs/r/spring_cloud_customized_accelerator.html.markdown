---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_customized_accelerator"
description: |-
  Manages a Spring Cloud Customized Accelerator.
---

# azurerm_spring_cloud_customized_accelerator

Manages a Spring Cloud Customized Accelerator.

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_customized_accelerator` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "west europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_accelerator" "example" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
}

resource "azurerm_spring_cloud_customized_accelerator" "example" {
  name                        = "example"
  spring_cloud_accelerator_id = azurerm_spring_cloud_accelerator.example.id

  git_repository {
    url                 = "https://github.com/Azure-Samples/piggymetrics"
    git_tag             = "spring.version.2.0.3"
    interval_in_seconds = 100
  }

  accelerator_tags = ["tag-a", "tag-b"]
  description      = "example description"
  display_name     = "example name"
  icon_url         = "https://images.freecreatives.com/wp-content/uploads/2015/05/smiley-559124_640.jpg"
}
```

## Arguments Reference

The following arguments are supported:

* `git_repository` - (Required) A `git_repository` block as defined below.

* `name` - (Required) The name which should be used for this Spring Cloud Customized Accelerator. Changing this forces a new Spring Cloud Customized Accelerator to be created.

* `spring_cloud_accelerator_id` - (Required) The ID of the Spring Cloud Accelerator. Changing this forces a new Spring Cloud Customized Accelerator to be created.

---

* `accelerator_tags` - (Optional) Specifies a list of accelerator tags.

* `accelerator_type` - (Optional) Specifies the type of the Spring Cloud Customized Accelerator. Possible values are `Accelerator` and `Fragment`. Defaults to `Accelerator`.

* `description` - (Optional) Specifies the description of the Spring Cloud Customized Accelerator.

* `display_name` - (Optional) Specifies the display name of the Spring Cloud Customized Accelerator..

* `icon_url` - (Optional) Specifies the icon URL of the Spring Cloud Customized Accelerator..

---

A `git_repository` block supports the following:

* `url` - (Required) Specifies Git repository URL for the accelerator.

* `basic_auth` - (Optional) A `basic_auth` block as defined below. Conflicts with `git_repository[0].ssh_auth`. Changing this forces a new Spring Cloud Customized Accelerator to be created.

* `branch` - (Optional) Specifies the Git repository branch to be used.

* `ca_certificate_id` - (Optional) Specifies the ID of the CA Spring Cloud Certificate for https URL of Git repository.

* `commit` - (Optional) Specifies the Git repository commit to be used.

* `git_tag` - (Optional) Specifies the Git repository tag to be used.

* `interval_in_seconds` - (Optional) Specifies the interval for checking for updates to Git or image repository. It should be greater than 10.

* `ssh_auth` - (Optional) A `ssh_auth` block as defined below. Conflicts with `git_repository[0].basic_auth`. Changing this forces a new Spring Cloud Customized Accelerator to be created.

* `path` - (Optional) Specifies the path under the git repository to be treated as the root directory of the accelerator or the fragment (depending on `accelerator_type`).

---

A `basic_auth` block supports the following:

* `password` - (Required) Specifies the password of git repository basic auth.

* `username` - (Required) Specifies the username of git repository basic auth.

---

A `ssh_auth` block supports the following:

* `private_key` - (Required) Specifies the Private SSH Key of git repository basic auth.

* `host_key` - (Optional) Specifies the Public SSH Key of git repository basic auth.

* `host_key_algorithm` - (Optional) Specifies the SSH Key algorithm of git repository basic auth.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Spring Cloud Customized Accelerator.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Customized Accelerator.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Customized Accelerator.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Customized Accelerator.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Customized Accelerator.

## Import

Spring Cloud Customized Accelerators can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_customized_accelerator.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/applicationAccelerators/default/customizedAccelerators/customizedAccelerator1
```
