---
subcategory: "Template"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_deployment_script_azure_cli"
description: |-
  Manages a Resource Deployment Script of Azure Cli.
---

# azurerm_resource_deployment_script_azure_cli

Manages a Resource Deployment Script of Azure Cli.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-uai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_resource_deployment_script_azure_cli" "example" {
  name                = "example-rdsac"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  version             = "2.40.0"
  retention_interval  = "P1D"
  command_line        = "'foo' 'bar'"
  cleanup_preference  = "OnSuccess"
  force_update_tag    = "1"
  timeout             = "PT30M"

  script_content = <<EOF
            echo "{\"name\":{\"displayName\":\"$1 $2\"}}" > $AZ_SCRIPTS_OUTPUT_PATH
  EOF

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }

  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Resource Deployment Script. The name length must be from 1 to 260 characters. The name can only contain alphanumeric, underscore, parentheses, hyphen and period, and it cannot end with a period. Changing this forces a new Resource Deployment Script to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Resource Deployment Script should exist. Changing this forces a new Resource Deployment Script to be created.

* `location` - (Required) Specifies the Azure Region where the Resource Deployment Script should exist. Changing this forces a new Resource Deployment Script to be created.

* `version` - (Required) Specifies the version of the Azure CLI that should be used in the format `X.Y.Z` (e.g. `2.30.0`). A canonical list of versions [is available from the Microsoft Container Registry API](https://mcr.microsoft.com/v2/azure-cli/tags/list). Changing this forces a new Resource Deployment Script to be created.

* `retention_interval` - (Required) Interval for which the service retains the script resource after it reaches a terminal state. Resource will be deleted when this duration expires. The time duration should be between `1` hour and `26` hours (inclusive) and should be specified in ISO 8601 format. Changing this forces a new Resource Deployment Script to be created.

* `command_line` - (Optional) Command line arguments to pass to the script. Changing this forces a new Resource Deployment Script to be created.

* `cleanup_preference` - (Optional) Specifies the cleanup preference when the script execution gets in a terminal state. Possible values are `Always`, `OnExpiration`, `OnSuccess`. Defaults to `Always`. Changing this forces a new Resource Deployment Script to be created.

* `container` - (Optional) A `container` block as defined below. Changing this forces a new Resource Deployment Script to be created.

* `environment_variable` - (Optional) An `environment_variable` block as defined below. Changing this forces a new Resource Deployment Script to be created.

* `force_update_tag` - (Optional) Gets or sets how the deployment script should be forced to execute even if the script resource has not changed. Can be current time stamp or a GUID. Changing this forces a new Resource Deployment Script to be created.

* `identity` - (Optional) An `identity` block as defined below. Changing this forces a new Resource Deployment Script to be created.

* `primary_script_uri` - (Optional) Uri for the script. This is the entry point for the external script. Changing this forces a new Resource Deployment Script to be created.

* `script_content` - (Optional) Script body. Changing this forces a new Resource Deployment Script to be created.

* `storage_account` - (Optional) A `storage_account` block as defined below. Changing this forces a new Resource Deployment Script to be created.

* `supporting_script_uris` - (Optional) Supporting files for the external script. Changing this forces a new Resource Deployment Script to be created.

* `timeout` - (Optional) Maximum allowed script execution time specified in ISO 8601 format. Needs to be greater than 0 and smaller than 1 day. Defaults to `P1D`. Changing this forces a new Resource Deployment Script to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Deployment Script.

---

A `container` block supports the following:

* `container_group_name` - (Optional) Container group name, if not specified then the name will get auto-generated. For more information, please refer to the [Container Configuration](https://learn.microsoft.com/en-us/rest/api/resources/deployment-scripts/create?tabs=HTTP#containerconfiguration) documentation.

---

An `environment_variable` block supports the following:

* `name` - (Required) Specifies the name of the environment variable.

* `secure_value` - (Optional) Specifies the value of the secure environment variable.

* `value` - (Optional) Specifies the value of the environment variable.

---

An `identity` block supports the following:

* `type` - (Required) Type of the managed identity. The only possible value is `UserAssigned`. Changing this forces a new resource to be created.

* `identity_ids` - (Required) Specifies the list of user-assigned managed identity IDs associated with the resource. Changing this forces a new resource to be created.

---

A `storage_account` block supports the following:

* `key` - (Required) Specifies the storage account access key.

* `name` - (Required) Specifies the storage account name.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Deployment Script.

* `outputs` - List of script outputs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Deployment Script.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Deployment Script.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Deployment Script.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Deployment Script.

## Import

Resource Deployment Script can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_deployment_script_azure_cli.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Resources/deploymentScripts/script1
```
