---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_linked_service_azure_databricks"
description: |-
  Manages a Linked Service (connection) between Azure Databricks and Azure Data Factory.
---

# azurerm_data_factory_linked_service_azure_databricks

Manages a Linked Service (connection) between Azure Databricks and Azure Data Factory.

## Example Usage with managed identity & new cluster

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "East US"
}

#Create a Linked Service using managed identity and new cluster config
resource "azurerm_data_factory" "example" {
  name                = "TestDtaFactory92783401247"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  identity {
    type = "SystemAssigned"
  }
}

#Create a databricks instance
resource "azurerm_databricks_workspace" "example" {
  name                = "databricks-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
}

resource "azurerm_data_factory_linked_service_azure_databricks" "msi_linked" {
  name                = "ADBLinkedServiceViaMSI"
  data_factory_name   = azurerm_data_factory.example.name
  resource_group_name = azurerm_resource_group.example.name
  description         = "ADB Linked Service via MSI"
  adb_domain          = "https://${azurerm_databricks_workspace.example.workspace_url}"

  msi_work_space_resource_id = azurerm_databricks_workspace.example.id

  new_cluster_config {
    node_type             = "Standard_NC12"
    cluster_version       = "5.5.x-gpu-scala2.11"
    min_number_of_workers = 1
    max_number_of_workers = 5
    driver_node_type      = "Standard_NC12"
    log_destination       = "dbfs:/logs"

    custom_tags = {
      custom_tag1 = "sct_value_1"
      custom_tag2 = "sct_value_2"
    }

    spark_config = {
      config1 = "value1"
      config2 = "value2"
    }

    spark_environment_variables = {
      envVar1 = "value1"
      envVar2 = "value2"
    }

    init_scripts = ["init.sh", "init2.sh"]
  }
}


```

## Example Usage with access token & existing cluster

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "East US"
}

#Link to an existing cluster via access token
resource "azurerm_data_factory" "example" {
  name                = "TestDtaFactory92783401247"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

#Create a databricks instance
resource "azurerm_databricks_workspace" "example" {
  name                = "databricks-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"
}

resource "azurerm_data_factory_linked_service_azure_databricks" "at_linked" {
  name                = "ADBLinkedServiceViaAccessToken"
  data_factory_name   = azurerm_data_factory.example.name
  resource_group_name = azurerm_resource_group.example.name
  description         = "ADB Linked Service via Access Token"
  existing_cluster_id = "0308-201146-sly615"

  access_token = "SomeDatabricksAccessToken"
  adb_domain   = "https://${azurerm_databricks_workspace.example.workspace_url}"
}
```

## Arguments Reference

The following arguments are supported:

* `adb_domain` - (Required) The domain URL of the databricks instance.

* `data_factory_name` - (Required) The Data Factory name in which to associate the Linked Service with. Changing this forces a new resource.

* `name` - (Required) Specifies the name of the Data Factory Linked Service. Changing this forces a new resource to be created. Must be unique within a data factory. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory Linked Service. Changing this forces a new resource.

---

You must specify exactly one of the following authentication blocks:

* `access_token` - (Optional) Authenticate to ADB via an access token.

* `key_vault_password` - (Optional) Authenticate to ADB via Azure Key Vault Linked Service as defined in the `key_vault_password` block below.

* `msi_work_space_resource_id` - (Optional) Authenticate to ADB via managed service identity.

---

You must specify exactly one of the following modes for cluster integration:

* `existing_cluster_id` - (Optional) The cluster_id of an existing cluster within the linked ADB instance.

* `instance_pool` - (Optional) Leverages an instance pool within the linked ADB instance as defined by  `instance_pool` block below.

* `new_cluster_config` - (Optional) Creates new clusters within the linked ADB instance as defined in the  `new_cluster_config` block below.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Data Factory Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Linked Service.

* `description` - (Optional) The description for the Data Factory Linked Service.

* `integration_runtime_name` - (Optional) The integration runtime reference to associate with the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

---

A `key_vault_password` block supports the following:

* `linked_service_name` - (Required) Specifies the name of an existing Key Vault Data Factory Linked Service.

* `secret_name` - (Required) Specifies the secret name in Azure Key Vault that stores ADB access token.

---

A `new_cluster_config` block supports the following:

* `cluster_version` - (Required) Spark version of a the cluster.

* `node_type` - (Required) Node type for the new cluster.

* `custom_tags` - (Optional) Tags for the cluster resource.

* `driver_node_type` - (Optional) Driver node type for the cluster.

* `init_scripts` - (Optional) User defined initialization scripts for the cluster.

* `log_destination` - (Optional) Location to deliver Spark driver, worker, and event logs.

* `spark_config` - (Optional) User-specified Spark configuration variables key-value pairs.

* `spark_environment_variables` - (Optional) User-specified Spark environment variables key-value pairs.

---

A `instance_pool` block supports the following:

* `instance_pool_id` - (Required) Identifier of the instance pool within the linked ADB instance.

* `cluster_version` - (Required) Spark version of a the cluster.

* `min_number_of_workers` - (Optional) The minimum number of worker nodes. Defaults to 1.

* `max_number_of_workers` - (Optional) The max number of worker nodes. Set this value if you want to enable autoscaling between the `min_number_of_workers` and this value. Omit this value to use a fixed number of workers defined in the `min_number_of_workers` property.

---
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Factory Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Linked Service.

## Import

Data Factory Linked Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_linked_service_azure_databricks.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/linkedservices/example
```
