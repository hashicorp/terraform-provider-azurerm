---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_data_flow"
description: |-
  Manages a Data Flow inside an Azure Data Factory.
---

# azurerm_data_factory_data_flow

Manages a Data Flow inside an Azure Data Factory.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_data_factory_linked_custom_service" "example" {
  name                 = "linked_service"
  data_factory_id      = azurerm_data_factory.example.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.example.primary_connection_string}"
}
JSON
}

resource "azurerm_data_factory_dataset_json" "example1" {
  name                = "dataset1"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_custom_service.example.name

  azure_blob_storage_location {
    container = "container"
    path      = "foo/bar/"
    filename  = "foo.txt"
  }

  encoding = "UTF-8"
}

resource "azurerm_data_factory_dataset_json" "example2" {
  name                = "dataset2"
  data_factory_id     = azurerm_data_factory.example.id
  linked_service_name = azurerm_data_factory_linked_custom_service.example.name

  azure_blob_storage_location {
    container = "container"
    path      = "foo/bar/"
    filename  = "bar.txt"
  }

  encoding = "UTF-8"
}

resource "azurerm_data_factory_data_flow" "example" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  source {
    name = "source1"

    flowlet {
      name = azurerm_data_factory_flowlet_data_flow.example1.name
      parameters = {
        "Key1" = "value1"
      }
    }

    dataset {
      name = azurerm_data_factory_dataset_json.example1.name
    }
  }

  sink {
    name = "sink1"

    flowlet {
      name = azurerm_data_factory_flowlet_data_flow.example2.name
      parameters = {
        "Key1" = "value1"
      }
    }

    dataset {
      name = azurerm_data_factory_dataset_json.example2.name
    }
  }

  script = <<EOT
source(
  allowSchemaDrift: true, 
  validateSchema: false, 
  limit: 100, 
  ignoreNoFilesFound: false, 
  documentForm: 'documentPerLine') ~> source1 
source1 sink(
  allowSchemaDrift: true, 
  validateSchema: false, 
  skipDuplicateMapInputs: true, 
  skipDuplicateMapOutputs: true) ~> sink1
EOT
}

resource "azurerm_data_factory_flowlet_data_flow" "example1" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  source {
    name = "source1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.example.name
    }
  }

  sink {
    name = "sink1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.example.name
    }
  }

  script = <<EOT
source(
  allowSchemaDrift: true, 
  validateSchema: false, 
  limit: 100, 
  ignoreNoFilesFound: false, 
  documentForm: 'documentPerLine') ~> source1 
source1 sink(
  allowSchemaDrift: true, 
  validateSchema: false, 
  skipDuplicateMapInputs: true, 
  skipDuplicateMapOutputs: true) ~> sink1
EOT
}

resource "azurerm_data_factory_flowlet_data_flow" "example2" {
  name            = "example"
  data_factory_id = azurerm_data_factory.example.id

  source {
    name = "source1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.example.name
    }
  }

  sink {
    name = "sink1"

    linked_service {
      name = azurerm_data_factory_linked_custom_service.example.name
    }
  }

  script = <<EOT
source(
  allowSchemaDrift: true, 
  validateSchema: false, 
  limit: 100, 
  ignoreNoFilesFound: false, 
  documentForm: 'documentPerLine') ~> source1 
source1 sink(
  allowSchemaDrift: true, 
  validateSchema: false, 
  skipDuplicateMapInputs: true, 
  skipDuplicateMapOutputs: true) ~> sink1
EOT
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory Data Flow. Changing this forces a new resource to be created.

* `data_factory_id` - (Required) The ID of Data Factory in which to associate the Data Flow with. Changing this forces a new resource.

* `script` - (Optional) The script for the Data Factory Data Flow.

* `script_lines` - (Optional) The script lines for the Data Factory Data Flow.

* `source` - (Required) One or more `source` blocks as defined below.

* `sink` - (Required) One or more `sink` blocks as defined below.

* `annotations` - (Optional) List of tags that can be used for describing the Data Factory Data Flow.

* `description` - (Optional) The description for the Data Factory Data Flow.

* `folder` - (Optional) The folder that this Data Flow is in. If not specified, the Data Flow will appear at the root level.

* `transformation` - (Optional) One or more `transformation` blocks as defined below.

---

A `source` block supports the following:

* `name` - (Required) The name for the Data Flow Source.

* `description` - (Optional) The description for the Data Flow Source.

* `dataset` - (Optional) A `dataset` block as defined below.

* `flowlet` - (Optional) A `flowlet` block as defined below.

* `linked_service` - (Optional) A `linked_service` block as defined below.

* `rejected_linked_service` - (Optional) A `rejected_linked_service` block as defined below.

* `schema_linked_service` - (Optional) A `schema_linked_service` block as defined below.

---

A `sink` block supports the following:

* `name` - (Required) The name for the Data Flow Source.

* `description` - (Optional) The description for the Data Flow Source.

* `dataset` - (Optional) A `dataset` block as defined below.

* `flowlet` - (Optional) A `flowlet` block as defined below.

* `linked_service` - (Optional) A `linked_service` block as defined below.

* `rejected_linked_service` - (Optional) A `rejected_linked_service` block as defined below.

* `schema_linked_service` - (Optional) A `schema_linked_service` block as defined below.

---

A `dataset` block supports the following:

* `name` - (Required) The name for the Data Factory Dataset.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory dataset.

---

A `flowlet` block supports the following:

* `name` - (Required) The name for the Data Factory Flowlet.

* `dataset_parameters` - (Optional) Specifies the reference data flow parameters from dataset.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Flowlet.

---

A `linked_service` block supports the following:

* `name` - (Required) The name for the Data Factory Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

---

A `rejected_linked_service` block supports the following:

* `name` - (Required) The name for the Data Factory Linked Service with schema.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

---

A `schema_linked_service` block supports the following:

* `name` - (Required) The name for the Data Factory Linked Service with schema.

* `parameters` - (Optional) A map of parameters to associate with the Data Factory Linked Service.

---

A `transformation` block supports the following:

* `name` - (Required) The name for the Data Flow transformation.

* `description` - (Optional) The description for the Data Flow transformation.

* `dataset` - (Optional) A `dataset` block as defined below.

* `flowlet` - (Optional) A `flowlet` block as defined below.

* `linked_service` - (Optional) A `linked_service` block as defined below.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory Data Flow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory Data Flow.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory Data Flow.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory Data Flow.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory Data Flow.

## Import

Data Factory Data Flow can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_data_flow.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example/dataflows/example
```
