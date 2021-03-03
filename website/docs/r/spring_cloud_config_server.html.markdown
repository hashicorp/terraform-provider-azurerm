---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_config_server"
description: |-
  Manages an Azure Spring Cloud Config Server.
---

# azurerm_spring_cloud_config_server

Manages an Azure Spring Cloud Config Server.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Southeast Asia"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "S0"
}

resource "azurerm_spring_cloud_config_server" "example" {
  spring_cloud_service_id = azurerm_spring_cloud_service.example.id
  uri                     = "https://github.com/Azure-Samples/piggymetrics"
  label                   = "config"
  search_paths            = ["dir1", "dir2"]
}
```

## Argument Reference

The following arguments are supported:

* `spring_cloud_service_id` - (Required) Specifies the id of the Spring Cloud Service resource in which to create the Spring Cloud Config Server. Changing this forces a new resource to be created.

* `uri` - (Required) The URI of the default Git repository used as the Config Server back end, should be started with `http://`, `https://`, `git@`, or `ssh://`.

* `http_basic_auth` - (Optional) A `http_basic_auth` block as defined below.

* `label` - (Optional) The default label of the Git repository, should be the branch name, tag name, or commit-id of the repository.

* `repository` - (Optional) One or more `repository` blocks as defined below.

* `search_paths` - (Optional) An array of strings used to search subdirectories of the Git repository.

* `ssh_auth` - (Optional) A `ssh_auth` block as defined below.

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

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Spring Cloud Config Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Config Server.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Config Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Config Server.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Config Server.

## Import

Spring Cloud Config Server can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_config_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AppPlatform/Spring/spring1
```
