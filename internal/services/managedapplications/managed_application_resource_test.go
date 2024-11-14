// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedapplications_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ManagedApplicationResource struct{}

func TestAccManagedApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("outputs.boolOutput").HasValue("true"),
				check.That(data.ResourceName).Key("outputs.intOutput").HasValue("100"),
				check.That(data.ResourceName).Key("outputs.objectOutput").HasValue("{\"nested_array\":[\"value_1\",\"value_2\"],\"nested_bool\":true,\"nested_object\":{\"key_0\":0}}"),
				check.That(data.ResourceName).Key("outputs.stringOutput").HasValue("stringOutputValue"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccManagedApplication_switchBetweenParametersAndParameterValues(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("skipping bacause `parameters` is deprecated in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.parameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_parameters(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("skipping as `parameters` is deprecated in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.parameters(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.parametersUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_parameterValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithParameterValuesUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_allSupportedParameterValuesTypes(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.allSupportedParameterValuesTypes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_parametersSecureString(t *testing.T) {
	if features.FourPointOhBeta() {
		t.Skipf("skipping because `parameters` is removed in 4.0")
	}

	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.parametersSecureString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parameters.secureStringParameter", "parameter_values"),
	})
}

func TestAccManagedApplication_parameterValuesSecureString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.parameterValuesSecureString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parameters.secureStringParameter", "parameter_values"),
	})
}

func TestAccManagedApplication_plan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}
	publisher := "cisco"
	offer := "cisco-meraki-vmx"
	plan := "cisco-meraki-vmx"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				if err := r.cancelExistingAgreement(t, publisher, offer, plan); err != nil {
					t.Fatalf("Failed to cancel existing agreement with error: %+v", err)
				}
			},
			Config: r.plan(data, publisher, offer, plan),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagedApplication_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_managed_application", "test")
	r := ManagedApplicationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ManagedApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := applications.ParseApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ManagedApplication.ApplicationClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ManagedApplicationResource) parameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameters = {
    stringParameter       = "value_1"
    secureStringParameter = ""
  }
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) parametersUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameters = {
    stringParameter       = "value_2"
    secureStringParameter = ""
  }
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_managed_application" "import" {
  name                        = azurerm_managed_application.test.name
  location                    = azurerm_managed_application.test.location
  resource_group_name         = azurerm_managed_application.test.resource_group_name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%d"
}
`, r.basic(data), data.RandomInteger)
}

func (r ManagedApplicationResource) plan(data acceptance.TestData, publisher string, offer string, plan string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mapp-%[2]d"
  location = "%[1]s"
}

resource "azurerm_marketplace_agreement" "test" {
  publisher = "%[3]s"
  offer     = "%[4]s"
  plan      = "%[5]s"
}

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "MarketPlace"
  managed_resource_group_name = "infraGroup%[2]d"

  plan {
    name      = azurerm_marketplace_agreement.test.plan
    product   = azurerm_marketplace_agreement.test.offer
    publisher = azurerm_marketplace_agreement.test.publisher
    version   = "15.37.1"
  }

  parameter_values = jsonencode({
    zone = {
      value = "0"
    },
    location = {
      value = azurerm_resource_group.test.location
    },
    merakiAuthToken = {
      value = "f451adfb-d00b-4612-8799-b29294217d4a"
    },
    subnetAddressPrefix = {
      value = "10.0.0.0/24"
    },
    subnetName = {
      value = "acctestSubnet"
    },
    virtualMachineSize = {
      value = "Standard_DS12_v2"
    },
    virtualNetworkAddressPrefix = {
      value = "10.0.0.0/16"
    },
    virtualNetworkName = {
      value = "acctestVnet"
    },
    virtualNetworkNewOrExisting = {
      value = "new"
    },
    virtualNetworkResourceGroup = {
      value = "acctestVnetRg"
    },
    vmName = {
      value = "acctestVM"
    }
  })
}
`, data.Locations.Primary, data.RandomInteger, publisher, offer, plan)
}

func (r ManagedApplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameter_values = jsonencode({
    stringParameter = {
      value = "value_1_from_parameter_values"
    },
    secureStringParameter = {
      value = ""
    }
  })
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) basicWithParameterValuesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameter_values = jsonencode({
    stringParameter = {
      value = "value_1_from_parameter_values_updated"
    },
    secureStringParameter = {
      value = ""
    }
  })
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) allSupportedParameterValuesTypes(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

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
      value = ""
    },
    objectParameter = {
      value = {
        nested_bool  = true
        nested_array = ["value_1", "value_2"]
        nested_object = {
          key_0 = 0
        }
      }
    },
    arrayParameter = {
      value = ["value_1", "value_2"]
    }
  })
}
`, r.templateAllSupportedParametersTypes(data), data.RandomInteger)
}

func (r ManagedApplicationResource) parametersSecureString(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameters = {
    stringParameter       = "value_1"
    secureStringParameter = "secure_value_1"
  }
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) parameterValuesSecureString(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

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
    },
    arrayParameter = {
      value = ["value_1", "value_2"]
    }
  })
}
`, r.templateAllSupportedParametersTypes(data), data.RandomInteger)
}

func (r ManagedApplicationResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameter_values = jsonencode({
    stringParameter = {
      value = "value_1_from_parameter_values"
    },
    secureStringParameter = {
      value = ""
    }
  })

  tags = {
    ENV = "Test"
  }
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) tagsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_managed_application" "test" {
  name                        = "acctestManagedApp%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  kind                        = "ServiceCatalog"
  managed_resource_group_name = "infraGroup%[2]d"
  application_definition_id   = azurerm_managed_application_definition.test.id

  parameter_values = jsonencode({
    stringParameter = {
      value = "value_1_from_parameter_values"
    },
    secureStringParameter = {
      value = ""
    }
  })

  tags = {
    ENV = "Test2"
  }
}
`, r.templateStringParameter(data), data.RandomInteger)
}

func (r ManagedApplicationResource) templateStringParameter(data acceptance.TestData) string {
	parameters := `
         "stringParameter": {
            "type": "string"
         },
         "secureStringParameter": {
            "type": "secureString"
         }
`
	return r.template(data, parameters)
}

func (r ManagedApplicationResource) templateAllSupportedParametersTypes(data acceptance.TestData) string {
	parameters := `
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
         },
         "arrayParameter": {
            "type": "array"
         }
`
	return r.template(data, parameters)
}

func (ManagedApplicationResource) template(data acceptance.TestData, parameters string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

data "azurerm_role_definition" "test" {
  name = "Contributor"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mapp-%[2]d"
  location = "%[1]s"
}

resource "azurerm_managed_application_definition" "test" {
  name                = "acctestManagedAppDef%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  lock_level          = "ReadOnly"
  display_name        = "TestManagedAppDefinition"
  description         = "Test Managed App Definition"
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
%[3]s
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
        },
        "arrayOutput": {
          "type": "array",
          "value": ["value_1", "value_2"]
        }
      }
    }
  MAIN_TEMPLATE

  authorization {
    service_principal_id = data.azurerm_client_config.test.object_id
    role_definition_id   = split("/", data.azurerm_role_definition.test.id)[length(split("/", data.azurerm_role_definition.test.id)) - 1]
  }
}
`, data.Locations.Primary, data.RandomInteger, parameters)
}

func (ManagedApplicationResource) cancelExistingAgreement(t *testing.T, publisher string, offer string, plan string) error {
	clientManager, err := testclient.Build()
	if err != nil {
		t.Fatalf("building client: %+v", err)
	}

	ctx, cancel := context.WithDeadline(clientManager.StopContext, time.Now().Add(15*time.Minute))
	defer cancel()

	client := clientManager.Compute.MarketplaceAgreementsClient
	subscriptionId := clientManager.Account.SubscriptionId

	idGet := agreements.NewOfferPlanID(subscriptionId, publisher, offer, plan)
	idCancel := agreements.NewPlanID(subscriptionId, publisher, offer, plan)

	existing, err := client.MarketplaceAgreementsGet(ctx, idGet)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", idGet, err)
	}

	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			if accepted := props.Accepted; accepted != nil && *accepted {
				resp, err := client.MarketplaceAgreementsCancel(ctx, idCancel)
				if err != nil {
					if response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("marketplace agreement %q does not exist", idGet)
					}
					return fmt.Errorf("canceling %s: %+v", idGet, err)
				}
			}
		}
	}

	return nil
}
