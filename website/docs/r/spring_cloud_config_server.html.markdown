subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_config_server"
sidebar_current: "docs-azurerm-resource-spring-cloud-config-server"
description: |-
  Manage Azure Spring Cloud Config Server instance.
---

# azurerm_spring_cloud

Manage Azure Cloud Config Server instance. Spring Cloud Config provides server and client-side support for an externalized configuration in a distributed system. With the Config Server instance, you have a central place to manage external properties for applications across all environments.

Azure Spring Cloud Config Server supports Azure DevOps, GitHub, GitLab, and Bitbucket for storing your Config Server files

Azure Spring Cloud Config Server supports Public repository, Private repository with SSH authentication or Private repository with basic authentication

Azure Spring Cloud Config Server also supports multiple repositories with Pattern Matching


## Argument Reference

The following arguments are supported:

* `spring_cloud_id` - (Required) The id of the Spring Cloud Service resource. Changing this forces a new resource to be created.

* `repositories` - (Optional) One or more `repository` block defined below.

* `uri` - (Required) The URI of the default Git repository used as the Config Server back end, should be started with http://, https://, git@, or ssh://.

* `label` - (Optional) The default label of the Git repository, should be the branch name, tag name, or commit-id of the repository.

* `search_paths` - (Optional) An array of strings used to search subdirectories of the Git repository.

* `username` - (Optional) The username that's used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

* `password` - (Optional) The password used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

* `host_key` - (Optional) The host key of the Git repository server, should not include the algorithm prefix as covered by host-key-algorithm.

* `host_key_algorithm` - (Optional) The host key algorithm, should be ssh-dss, ssh-rsa, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, or ecdsa-sha2-nistp521. Required only if host-key exists.

* `private_key` - (Optional) The SSH private key to access the Git repository, required when the URI starts with git@ or ssh://.

* `strict_host_key_checking` - (Optional) Indicates whether the Config Server instance will fail to start when leveraging the private host-key. Should be true (default value) or false.


---

The `repository` block supports the following:

* `name` - (Required) A name to identify on the Git repository, required only if repos exists. For example, team-A, team-B

* `pattern` - (Optional) An array of strings used to match an application name. For each pattern, use the {application}/{profile} format with wildcards.

* `uri` - (Required) The URI of the Git repository that's used as the Config Server back end should be started with http://, https://, git@, or ssh://.

* `label` - (Optional) The default label of the Git repository, should be the branch name, tag name, or commit-id of the repository.

* `search_paths` - (Optional) An array of strings used to search subdirectories of the Git repository.

* `username` - (Optional) The username that's used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

* `password` - (Optional) The password used to access the Git repository server, required when the Git repository server supports Http Basic Authentication.

* `host_key` - (Optional) The host key of the Git repository server, should not include the algorithm prefix as covered by host-key-algorithm.

* `host_key_algorithm` - (Optional) The host key algorithm, should be ssh-dss, ssh-rsa, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, or ecdsa-sha2-nistp521. Required only if host-key exists.

* `private_key` - (Optional) The SSH private key to access the Git repository, required when the URI starts with git@ or ssh://.

* `strict_host_key_checking` - (Optional) Indicates whether the Config Server instance will fail to start when leveraging the private host-key. Should be true (default value) or false.

## Attributes Reference

The following attributes are exported:

* `id` - Fully qualified resource Id for the resource.
