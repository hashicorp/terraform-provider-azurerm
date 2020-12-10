package managedapplications_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managedapplications/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ManagedApplicationDefinitionResource struct {
}

func TestAccManagedApplicationDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application_definition", "test")
	r := ManagedApplicationDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("package_file_uri"),
	})
}

func TestAccManagedApplicationDefinition_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application_definition", "test")
	r := ManagedApplicationDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccManagedApplicationDefinition_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application_definition", "test")
	r := ManagedApplicationDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("create_ui_definition", "main_template"),
	})
}

func TestAccManagedApplicationDefinition_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application_definition", "test")
	r := ManagedApplicationDefinitionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("TestManagedApplicationDefinition"),
				check.That(data.ResourceName).Key("description").HasValue("Test Managed Application Definition"),
				check.That(data.ResourceName).Key("package_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("package_file_uri").Exists(),
			),
		},
		data.ImportStep("package_file_uri"),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("UpdatedTestManagedApplicationDefinition"),
				check.That(data.ResourceName).Key("description").HasValue("Updated Test Managed Application Definition"),
				check.That(data.ResourceName).Key("package_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
				check.That(data.ResourceName).Key("create_ui_definition").Exists(),
				check.That(data.ResourceName).Key("main_template").Exists(),
			),
		},
		data.ImportStep("create_ui_definition", "main_template", "package_file_uri"),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("TestManagedApplicationDefinition"),
				check.That(data.ResourceName).Key("description").HasValue("Test Managed Application Definition"),
				check.That(data.ResourceName).Key("package_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("package_file_uri").Exists(),
			),
		},
		data.ImportStep("create_ui_definition", "main_template", "package_file_uri"),
	})
}

func (ManagedApplicationDefinitionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ManagedApplication.ApplicationDefinitionClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Managed Definition %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ApplicationDefinitionProperties != nil), nil
}

func (r ManagedApplicationDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application_definition" "test" {
  name                = "acctestAppDef%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lock_level          = "None"
  package_file_uri    = "https://github.com/Azure/azure-managedapp-samples/raw/master/Managed Application Sample Packages/201-managed-storage-account/managedstorage.zip"
  display_name        = "TestManagedApplicationDefinition"
  description         = "Test Managed Application Definition"
  package_enabled     = false
}
`, r.template(data), data.RandomInteger)
}

func (r ManagedApplicationDefinitionResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application_definition" "import" {
  name                = azurerm_managed_application_definition.test.name
  location            = azurerm_managed_application_definition.test.location
  resource_group_name = azurerm_managed_application_definition.test.resource_group_name
  display_name        = azurerm_managed_application_definition.test.display_name
  lock_level          = azurerm_managed_application_definition.test.lock_level
}
`, r.basic(data))
}

func (r ManagedApplicationDefinitionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application_definition" "test" {
  name                = "acctestAppDef%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lock_level          = "ReadOnly"
  display_name        = "UpdatedTestManagedApplicationDefinition"
  description         = "Updated Test Managed Application Definition"
  package_enabled     = true

  create_ui_definition = <<CREATE_UI_DEFINITION
    {
      "$schema": "https://schema.management.azure.com/schemas/0.1.2-preview/CreateUIDefinition.MultiVm.json#",
      "handler": "Microsoft.Azure.CreateUIDef",
	  "version": "0.1.2-preview",
	  "parameters": {
         "basics": [
            {}
         ],
         "steps": [
            {
              "name": "storageConfig",
              "label": "Storage settings",
              "subLabel": {
                 "preValidation": "Configure the infrastructure settings",
                 "postValidation": "Done"
              },
              "bladeTitle": "Storage settings",
              "elements": [
                 {
                   "name": "storageAccounts",
                   "type": "Microsoft.Storage.MultiStorageAccountCombo",
                   "label": {
                      "prefix": "Storage account name prefix",
                      "type": "Storage account type"
                   },
                   "defaultValue": {
                      "type": "Standard_LRS"
                   },
                   "constraints": {
                      "allowedTypes": [
                        "Premium_LRS",
                        "Standard_LRS",
                        "Standard_GRS"
                      ]
                   }
                 }
              ]
            }
         ],
         "outputs": {
            "storageAccountNamePrefix": "[steps('storageConfig').storageAccounts.prefix]",
            "storageAccountType": "[steps('storageConfig').storageAccounts.type]",
            "location": "[location()]"
         }
      }
	}
  CREATE_UI_DEFINITION

  main_template = <<MAIN_TEMPLATE
    {
      "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
      "contentVersion": "1.0.0.0",
      "parameters": {
         "storageAccountNamePrefix": {
            "type": "string"
         },
         "storageAccountType": {
            "type": "string"
         },
         "location": {
            "type": "string",
            "defaultValue": "[resourceGroup().location]"
         }
      },
      "variables": {
         "storageAccountName": "[concat(parameters('storageAccountNamePrefix'), uniqueString(resourceGroup().id))]"
      },
      "resources": [
         {
           "type": "Microsoft.Storage/storageAccounts",
           "name": "[variables('storageAccountName')]",
           "apiVersion": "2016-01-01",
           "location": "[parameters('location')]",
           "sku": {
              "name": "[parameters('storageAccountType')]"
           },
           "kind": "Storage",
           "properties": {}
         }
      ],
      "outputs": {
         "storageEndpoint": {
           "type": "string",
           "value": "[reference(resourceId('Microsoft.Storage/storageAccounts/', variables('storageAccountName')), '2016-01-01').primaryEndpoints.blob]"
         }
      }
    }
  MAIN_TEMPLATE

  authorization {
    service_principal_id = data.azurerm_client_config.current.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.builtin.id)[length(split("/", data.azurerm_role_definition.builtin.id)) - 1]
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (ManagedApplicationDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_role_definition" "builtin" {
  name = "Contributor"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mapp-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
