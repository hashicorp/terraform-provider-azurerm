---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_service"
description: |-
  Manages an Azure Spring Cloud Service.
---

# azurerm_spring_cloud_service

Manages an Azure Spring Cloud Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "S0"

  config_server_git_setting {
    uri          = "https://github.com/Azure-Samples/piggymetrics"
    label        = "config"
    search_paths = ["dir1", "dir2"]
  }

  trace {
    instrumentation_key = azurerm_application_insights.example.instrumentation_key
    sample_rate         = 10.0
  }

  tags = {
    Env = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies The name of the resource group in which to create the Spring Cloud Service. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

-> **Note:** At this time Azure Spring Cloud Service is only supported in a subset of regions (including `East US`, `South East Asia`, `West Europe` and `West US 2`.

* `sku_name` - (Optional) Specifies the SKU Name for this Spring Cloud Service. Possible values are `B0` and `S0`. Defaults to `S0`.

* `network` - (Optional) A `network` block as defined below. Changing this forces a new resource to be created.

* `config_server_git_setting` - (Optional) A `config_server_git_setting` block as defined below.

* `trace` - (Optional) A `trace` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `network` block supports the following:

* `app_subnet_id` - (Required) Specifies the ID of the Subnet which should host the Spring Boot Applications deployed in this Spring Cloud Service. Changing this forces a new resource to be created.

* `service_runtime_subnet_id` - (Required) Specifies the ID of the Subnet where the Service Runtime components of the Spring Cloud Service will exist. Changing this forces a new resource to be created.

* `cidr_ranges` - (Required) A list of (at least 3) CIDR ranges (at least /16) which are used to host the Spring Cloud infrastructure, which must not overlap with any existing CIDR ranges in the Subnet. Changing this forces a new resource to be created.

* `app_network_resource_group` - (Optional) Specifies the Name of the resource group containing network resources of Azure Spring Cloud Apps. Changing this forces a new resource to be created.

* `service_runtime_network_resource_group` - (Optional) Specifies the Name of the resource group containing network resources of Azure Spring Cloud Service Runtime. Changing this forces a new resource to be created.

---

The `config_server_git_setting` block supports the following:

* `uri` - (Required) The URI of the default Git repository used as the Config Server back end, should be started with `http://`, `https://`, `git@`, or `ssh://`.

* `label` - (Optional) The default label of the Git repository, should be the branch name, tag name, or commit-id of the repository.

* `search_paths` - (Optional) An array of strings used to search subdirectories of the Git repository.

* `http_basic_auth` - (Optional) A `http_basic_auth` block as defined below.

* `ssh_auth` - (Optional) A `ssh_auth` block as defined below.

* `repository` - (Optional) One or more `repository` blocks as defined below.

---

The `repository` block supports the following:

* `name` - (Required) A name to identify on the Git repository, required only if repos exists.

* `uri` - (Required) The URI of the Git repository that's used as the Config Server back end should be started with `http://`, `https://`, `git@`, or `ssh://`.

* `pattern` - (Optional) An array of strings used to match an application name. For each pattern, use the `{application}/{profile}` format with wildcards.

* `label` - (Optional) The default label of the Git repository, should be the branch name, tag name, or commit-id of the repository.

* `search_paths` - (Optional) An array of strings used to search subdirectories of the Git repository.

* `http_basic_auth` - (Optional) A `http_basic_auth` block as defined below.

* `ssh_auth` - (Optional) A `ssh_auth` block as defined below.

---

The `http_basic_auth` block supports the following:

* `username` - (Required) The username that's used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

* `password` - (Required) The password used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

---

The `ssh_auth` block supports the following:

* `private_key` - (Required) The SSH private key to access the Git repository, required when the URI starts with `git@` or `ssh://`.

* `host_key` - (Optional) The host key of the Git repository server, should not include the algorithm prefix as covered by `host-key-algorithm`.

* `host_key_algorithm` - (Optional) The host key algorithm, should be `ssh-dss`, `ssh-rsa`, `ecdsa-sha2-nistp256`, `ecdsa-sha2-nistp384`, or `ecdsa-sha2-nistp521`. Required only if `host-key` exists.

* `strict_host_key_checking_enabled` - (Optional) Indicates whether the Config Server instance will fail to start if the host_key does not match.

---

The `trace` block supports the following:

* `instrumentation_key` - (Required) The Instrumentation Key used for Application Insights.

* `sample_rate` - (Optional) The sampling rate of Application Insights Agent. Must be between `0.0` and `100.0`. Defaults to `10.0`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Service.

* `outbound_public_ip_addresses` - A list of the outbound Public IP Addresses used by this Spring Cloud Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Spring Cloud Service.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Service.

## Import

Spring Cloud services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AppPlatform/Spring/spring1
```
