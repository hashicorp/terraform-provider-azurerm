# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_managed_application_definition" "example" {
  name                = "${var.prefix}managedappdef"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  lock_level          = "ReadOnly"
  display_name        = "ExampleManagedAppDefinition"
  description         = "Example Managed App Definition"
  package_enabled     = true

  create_ui_definition = <<CREATE_UI_DEFINITION
    {
      "$schema": "https://schema.management.azure.com/schemas/0.1.2-preview/CreateUIDefinition.MultiVm.json#",
      "handler": "Microsoft.Azure.CreateUIDef",
      "version": "0.1.2-preview",
      "parameters": {
         "basics": [],
         "steps": [],
         "outputs": {}
      }
    }
  CREATE_UI_DEFINITION

  main_template = <<MAIN_TEMPLATE
    {
      "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
      "contentVersion": "1.0.0.0",
      "parameters": {
         "boolParameter": {
            "type": "bool"
         },
         "intParameter": {
            "type": "int"
         },
         "stringParameter": {
            "type": "string"
         },
         "secureStringParameter": {
            "type": "secureString"
         },
         "objectParameter": {
            "type": "object"
         }
      },
      "variables": {},
      "resources": [],
      "outputs": {
        "boolOutput": {
          "type": "bool",
          "value": true
        },
        "intOutput": {
          "type": "int",
          "value": 100
        },
        "stringOutput": {
          "type": "string",
          "value": "stringOutputValue"
        },
        "objectOutput": {
          "type": "object",
          "value": {
            "nested_bool": true,
            "nested_array": ["value_1", "value_2"],
            "nested_object": {
              "key_0": 0
            }
          }
        }
      }
    }
  MAIN_TEMPLATE

  authorization {
    service_principal_id = data.azurerm_client_config.current.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.builtin.id)[length(split("/", data.azurerm_role_definition.builtin.id)) - 1]
  }
}


resource "azurerm_managed_application" "test" {
  name                        = "${var.prefix}managedapp"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "${var.prefix}infraGroup"
  application_definition_id   = azurerm_managed_application_definition.example.id

  parameter_values = jsonencode({
    boolParameter = {
      value = true
    },
    intParameter = {
      value = 100
    },
    stringParameter = {
      value = "value_1"
    },
    secureStringParameter = {
      value = "secure_value_1"
    },
    objectParameter = {
      value = {
        nested_bool  = true
        nested_array = ["value_1", "value_2"]
        nested_object = {
          key_0 = 0
        }
      }
    }
  })
}
