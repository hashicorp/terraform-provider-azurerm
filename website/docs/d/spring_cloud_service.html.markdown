---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_service"
description: |-
  Gets information about an existing Spring Cloud Service
---

# Data Source: azurerm_spring_cloud_service

Use this data source to access information about an existing Spring Cloud Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_spring_cloud_service" "example" {
  name                = azurerm_spring_cloud_service.example.name
  resource_group_name = azurerm_spring_cloud_service.example.resource_group_name
}

output "spring_cloud_service_id" {
  value = "${data.azurerm_spring_cloud_service.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies The name of the Spring Cloud Service resource.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Spring Cloud Service exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of Spring Cloud Service.

* `config_server_git_setting` - A `config_server_git_setting` block as defined below.

* `location` - The location of Spring Cloud Service.

* `outbound_public_ip_addresses` - A list of the outbound Public IP Addresses used by this Spring Cloud Service.

* `required_network_traffic_rules` - A list of `required_network_traffic_rules` blocks as defined below.

* `tags` - A mapping of tags assigned to Spring Cloud Service.

---

The `config_server_git_setting` block supports the following:

* `uri` - The URI of the Git repository

* `label` - The default label of the Git repository, which is a branch name, tag name, or commit-id of the repository

* `search_paths` - An array of strings used to search subdirectories of the Git repository.

* `http_basic_auth` - A `http_basic_auth` block as defined below.

* `ssh_auth` - A `ssh_auth` block as defined below.

* `repository` - One or more `repository` blocks as defined below.

---

The `repository` block contains the following:

* `name` - The name to identify on the Git repository.

* `pattern` - An array of strings used to match an application name. For each pattern, use the `{application}/{profile}` format with wildcards.

* `uri` - The URI of the Git repository

* `label` - The default label of the Git repository, which is a branch name, tag name, or commit-id of the repository

* `search_paths` - An array of strings used to search subdirectories of the Git repository.

* `http_basic_auth` - A `http_basic_auth` block as defined below.

* `ssh_auth` - A `ssh_auth` block as defined below.

---

The `http_basic_auth` block supports the following:

* `username` - The username used to access the Http Basic Authentication Git repository server.

* `password` - The password used to access the Http Basic Authentication Git repository server.

---

The `ssh_auth` block supports the following:

* `private_key` - The SSH private key to access the Git repository, needed when the URI starts with `git@` or `ssh://`.

* `host_key` - The host key of the Git repository server.

* `host_key_algorithm` - The host key algorithm.

* `strict_host_key_checking_enabled` - Indicates whether the Config Server instance will fail to start if the host_key does not match.

---

The `required_network_traffic_rules` supports the following:

* `direction` - The direction of required traffic. Possible values are `Inbound`, `Outbound`.

* `fqdns` - The FQDN list of required traffic.

* `ips` - The ip list of required traffic.

* `port` - The port of required traffic.

* `protocol` - The protocol of required traffic.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Service.
