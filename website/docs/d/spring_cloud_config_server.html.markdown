subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_config_server"
sidebar_current: "docs-azurerm-datasource-spring-cloud-config-server"
description: |-
  Gets information about an existing Spring Cloud Config Server
---

# Data Source: azurerm_spring_cloud

Use this data source to access information about an existing Spring Cloud Config Server.


## Argument Reference

The following arguments are supported:

* `spring_cloud_id` - (Required) The id of the Spring Cloud Service resource.

## Attributes Reference

The following attributes are exported:

* `repositories` - One or more `repository` block defined below.

* `uri` - URI of the repository

* `label` - Label of the repository

* `search_paths` - Searching path of the repository

* `username` - Username of git repository basic auth.

* `password` - Password of git repository basic auth.

* `host_key` - Public sshKey of git repository.

* `host_key_algorithm` - SshKey algorithm of git repository.

* `private_key` - Private sshKey algorithm of git repository.

* `strict_host_key_checking` - Strict host key checking or not.


---

The `repository` block contains the following:

* `name` - Name of the repository

* `pattern` - Collection of pattern of the repository

* `uri` - URI of the repository

* `label` - Label of the repository

* `search_paths` - Searching path of the repository

* `username` - Username of git repository basic auth.

* `password` - Password of git repository basic auth.

* `host_key` - Public sshKey of git repository.

* `host_key_algorithm` - SshKey algorithm of git repository.

* `private_key` - Private sshKey algorithm of git repository.

* `strict_host_key_checking` - Strict host key checking or not.
