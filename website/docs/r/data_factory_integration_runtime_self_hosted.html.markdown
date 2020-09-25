---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_integration_runtime_self_hosted"
description: |-
  Manages a Data Factory Self-hosted Integration Runtime.
---

# azurerm_data_factory_integration_runtime_self_hosted

Manages a Data Factory Self-hosted Integration Runtime.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "eastus"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_integration_runtime_self_hosted" "example" {
  name                = "example"
  resource_group_name = "example"
  data_factory_name   = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `data_factory_name` - (Required) Changing this forces a new Data Factory Self-hosted Integration Runtime to be created.

* `name` - (Required) The name which should be used for this Data Factory. Changing this forces a new Data Factory Self-hosted Integration Runtime to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Data Factory should exist. Changing this forces a new Data Factory Self-hosted Integration Runtime to be created.

---

* `description` - (Optional) Integration runtime description.

* `rbac_authorization` - (Optional) A `rbac_authorization` block as defined below.

---

A `rbac_authorization` block supports the following:

* `resource_id` - (Required) The resource identifier of the integration runtime to be shared. Changing this forces a new Data Factory to be created.

-> **Please Note**: RBAC Authorization creates a [linked Self-hosted Integration Runtime targeting the Shared Self-hosted Integration Runtime in resource_id](https://docs.microsoft.com/en-us/azure/data-factory/create-shared-self-hosted-integration-runtime-powershell#share-the-self-hosted-integration-runtime-with-another-data-factory). The linked Self-hosted Integration Runtime needs Contributor access granted to the Shared Self-hosted Data Factory. See example [Shared Self-hosted](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/data-factory/shared-self-hosted).

For more information on the configuration, please check out the [Azure documentation](https://docs.microsoft.com/en-us/rest/api/datafactory/integrationruntimes/createorupdate#linkedintegrationruntimerbacauthorization)

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Factory.

* `auth_key_1` - The primary integration runtime authentication key.

* `auth_key_2` - The secondary integration runtime authentication key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory.

## Import

Data Factories can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_integration_runtime_self_hosted.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/integrationruntimes/example
```
