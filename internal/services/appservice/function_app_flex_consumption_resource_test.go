// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FunctionAppFlexConsumptionResource struct{}

// remove in 5.0 starts
func TestAccFunctionAppFlexConsumption_FourPointOhBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAddBackendStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.addBackendStorageFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhBackendStorageUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.addBackendStorageFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.backendStorageUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhBackendStorageUseMsi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendStorageUseMsiFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhBackendStorageUseKeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendStorageUseKeyVaultFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhDeploymentStorageUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageBlobEndpointUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhDeploymentStorageSchemaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageSchemaUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhDeploymentStorageUaiSchemaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.deploymentStorageUai1NoBackendStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageUai2NoBackendStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageUseUaiWithBackendSaFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.storageUserAssignedIdentity2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhDeploymentStorageSystemIdentitySchemaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageSystemAssignedIdentityFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.storageSystemAssignedIdentityWithBackendSaFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.storageSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhConnectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStringFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhStickySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.stickySettingsFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhConnectionStringUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionStringFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStringUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStringFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAppSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAppSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettingsFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsAddKvpsFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsRemoveKvpsFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhRuntimePython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.pythonFourPointOh(data, "3.10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhInstanceMemoryUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.instanceMemoryUpdateFourPointOh(data, "20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAlwaysReadyUpdateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyBasicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCountFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAlwaysReadyInstanceCountError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.alwaysReadyInstanceCountErrorFourPointOh(data),
			ExpectError: regexp.MustCompile("the total number of always-ready instances should not exceed the maximum scale out limit"),
		},
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhAlwaysReadyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.alwaysReadyInstanceCountFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCountUpdateFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhMaxInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCountFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_FourPointOhHttpsOnlyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.completeFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basicFourPointOh(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

// remove in 5.0 ends

func TestAccFunctionAppFlexConsumption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_backendStorageUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.backendStorageUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_backendStorageUseMsi(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendStorageUseMsi(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_backendStorageUsingKeyVaultString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backendStorageUseKeyVault(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_deploymentStorageUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageBlobEndpointUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.deploymentStorageUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_connectionString(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_stickySettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.stickySettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_connectionStringUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionStringUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.connectionString(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_appSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_appSettingsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.appSettings(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsAddKvps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.appSettingsRemoveKvps(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimePython(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.python(data, "3.10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeNode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.runtimeNode(data, "20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_instanceMemoryUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.instanceMemoryUpdate(data, "20"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_alwaysReadyUpdateName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_alwaysReadyInstanceCountError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.alwaysReadyInstanceCountError(data),
			ExpectError: regexp.MustCompile("the total number of always-ready instances should not exceed the maximum scale out limit"),
		},
	})
}

func TestAccFunctionAppFlexConsumption_alwaysReadyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.alwaysReadyInstanceCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCountUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_maxInstanceCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.alwaysReadyInstanceCount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeJava(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeJavaUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.java(data, "17"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.java(data, "11"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimeDotNet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dotNet(data, "8.0"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_runtimePowerShell(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.powerShell(data, "7.4"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_systemAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageSystemAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_userAssignedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageUserAssignedIdentity1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_userAssignedIdentityUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageUserAssignedIdentity1(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.storageUserAssignedIdentity2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_httpsOnlyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kind").HasValue("functionapp,linux"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func TestAccFunctionAppFlexConsumption_vNetIntegrationWithVnetProperties(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_function_app_flex_consumption", "test")
	r := FunctionAppFlexConsumptionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.vNetIntegration_subnetWithVnetProperties(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("virtual_network_subnet_id").MatchesOtherKey(
					check.That("azurerm_subnet.test1").Key("id"),
				),
				check.That(data.ResourceName).Key("site_config.0.vnet_route_all_enabled").HasValue("true"),
			),
		},
		data.ImportStep("site_credential.0.password"),
	})
}

func (r FunctionAppFlexConsumptionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseFunctionAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.WebAppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

// remove in 5.0 starts
func (r FunctionAppFlexConsumptionResource) basicFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test1.primary_blob_endpoint}${azurerm_storage_container.test1-1.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test1.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageSchemaUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageBlobEndpointUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test1.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageUai1NoBackendStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctest-uai1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test1.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageUai2NoBackendStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageUseUaiWithBackendSaFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageSystemAssignedIdentityFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "SystemAssignedIdentity"
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageSystemAssignedIdentityWithBackendSaFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "SystemAssignedIdentity"
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) addBackendStorageFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test1.name
  storage_account_access_key = azurerm_storage_account.test1.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUseMsiFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name                  = azurerm_storage_account.test.name
  storage_account_uses_managed_identity = true

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key

  runtime_name           = "node"
  runtime_version        = "20"
  maximum_instance_count = 50
  instance_memory_in_mb  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUseKeyVaultFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  tags = {
    environment = "AccTest"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[2]s"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key};EndpointSuffix=core.windows.net"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  key_vault_reference_identity_id     = azurerm_user_assigned_identity.test.id
  storage_account_key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key

  runtime_name           = "node"
  runtime_version        = "20"
  maximum_instance_count = 50
  instance_memory_in_mb  = 2048

  site_config {}
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.templateKeyVault(data), data.RandomString, data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionStringFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) stickySettingsFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret"]
    connection_string_names = ["First", "Third"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionStringUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "AnotherExample"
    value = "some-other-connection-string"
    type  = "Custom"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) completeFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 100
  instance_memory_in_mb       = 2048
  https_only                  = true

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {
    app_command_line                       = "whoami"
    api_definition_url                     = "https://example.com/azure_function_app_def.json"
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]

    http2_enabled = true

    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    load_balancing_mode      = "LeastResponseTime"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"

    scm_ip_restriction {
      ip_address = "10.20.20.20/32"
      name       = "test-scm-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    scm_ip_restriction {
      ip_address = "fd80::/64"
      name       = "test-scm-restriction-v6"
      priority   = 124
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    websockets_enabled                = true
    health_check_path                 = "/health-check"
    health_check_eviction_time_in_min = 7
    worker_count                      = 3

    minimum_tls_version     = "1.2"
    scm_minimum_tls_version = "1.2"

    cors {
      allowed_origins = [
        "https://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsAddKvpsFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
    "tftestkvp2" : "tftestkvpvalue2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsRemoveKvpsFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) pythonFourPointOh(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "python"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, pythonVersion)
}

func (r FunctionAppFlexConsumptionResource) instanceMemoryUpdateFourPointOh(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "%s"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 4096

  site_config {}
}
`, r.template(data), data.RandomInteger, nodeVersion)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyBasicFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = azurerm_storage_container.test.id
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 100
  instance_memory_in_mb             = 2048
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCountFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = azurerm_storage_container.test.id
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 100
  instance_memory_in_mb             = 2048
  always_ready {
    name           = "blob"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCountErrorFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = azurerm_storage_container.test.id
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048
  always_ready {
    name           = "blob"
    instance_count = 20
  }
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 50
  }

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCountUpdateFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = azurerm_storage_container.test.id
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "20"
  maximum_instance_count            = 100
  instance_memory_in_mb             = 2048
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 20
  }
  always_ready {
    name           = "blob"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) vNetIntegration_subnetWithVnetPropertiesFourPointOh(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  lifecycle {
    ignore_changes = [
      delegation[0].service_delegation[0].actions
    ]
  }
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  virtual_network_subnet_id  = azurerm_subnet.test1.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 100
  instance_memory_in_mb       = 2048

  site_config {
    vnet_route_all_enabled = true
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

// remove in 5.0 ends

func (r FunctionAppFlexConsumptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageBlobEndpointUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test1.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) deploymentStorageUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test1.primary_blob_endpoint}${azurerm_storage_container.test1-1.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test1.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test1.name
  storage_account_access_key = azurerm_storage_account.test1.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUseMsi(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name                  = azurerm_storage_account.test.name
  storage_account_uses_managed_identity = true

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key

  runtime_name           = "node"
  runtime_version        = "20"
  maximum_instance_count = 50
  instance_memory_in_mb  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) backendStorageUseKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

%[1]s

data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%[2]s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get",
    ]

    secret_permissions = [
      "Get",
      "Delete",
      "List",
      "Purge",
      "Recover",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    secret_permissions = [
      "Get",
      "List",
    ]
  }

  tags = {
    environment = "AccTest"
  }
}

resource "azurerm_key_vault_secret" "test" {
  name         = "secret-%[2]s"
  value        = "DefaultEndpointsProtocol=https;AccountName=${azurerm_storage_account.test.name};AccountKey=${azurerm_storage_account.test.primary_access_key};EndpointSuffix=core.windows.net"
  key_vault_id = azurerm_key_vault.test.id
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-tf%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  key_vault_reference_identity_id     = azurerm_user_assigned_identity.test.id
  storage_account_key_vault_secret_id = azurerm_key_vault_secret.test.versionless_id

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key

  runtime_name           = "node"
  runtime_version        = "20"
  maximum_instance_count = 50
  instance_memory_in_mb  = 2048

  site_config {}
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.templateKeyVault(data), data.RandomString, data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionString(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) stickySettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  app_settings = {
    foo    = "bar"
    secret = "sauce"
    third  = "degree"
  }

  connection_string {
    name  = "First"
    value = "first-connection-string"
    type  = "Custom"
  }

  connection_string {
    name  = "Second"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "Third"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  sticky_settings {
    app_setting_names       = ["foo", "secret"]
    connection_string_names = ["First", "Third"]
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) connectionStringUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  connection_string {
    name  = "AnotherExample"
    value = "some-other-connection-string"
    type  = "Custom"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 100
  instance_memory_in_mb                  = 2048
  https_only                             = true

  app_settings = {
    foo    = "bar"
    secret = "sauce"
  }

  connection_string {
    name  = "Example"
    value = "some-postgresql-connection-string"
    type  = "PostgreSQL"
  }

  site_config {
    app_command_line                       = "whoami"
    api_definition_url                     = "https://example.com/azure_function_app_def.json"
    application_insights_connection_string = azurerm_application_insights.test.connection_string

    default_documents = [
      "first.html",
      "second.jsp",
      "third.aspx",
      "hostingstart.html",
    ]

    http2_enabled = true

    ip_restriction {
      ip_address = "10.10.10.10/32"
      name       = "test-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    load_balancing_mode      = "LeastResponseTime"
    remote_debugging_enabled = true
    remote_debugging_version = "VS2022"

    scm_ip_restriction {
      ip_address = "10.20.20.20/32"
      name       = "test-scm-restriction"
      priority   = 123
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    scm_ip_restriction {
      ip_address = "fd80::/64"
      name       = "test-scm-restriction-v6"
      priority   = 124
      action     = "Allow"
      headers {
        x_azure_fdid      = ["55ce4ed1-4b06-4bf1-b40e-4638452104da"]
        x_fd_health_probe = ["1"]
        x_forwarded_for   = ["9.9.9.9/32", "2002::1234:abcd:ffff:c0a8:101/64"]
        x_forwarded_host  = ["example.com"]
      }
    }

    websockets_enabled                = true
    health_check_path                 = "/health-check"
    health_check_eviction_time_in_min = 7
    worker_count                      = 3

    minimum_tls_version     = "1.2"
    scm_minimum_tls_version = "1.2"

    cors {
      allowed_origins = [
        "https://www.contoso.com",
        "www.contoso.com",
      ]

      support_credentials = true
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsAddKvps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
    "tftestkvp2" : "tftestkvpvalue2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) appSettingsRemoveKvps(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}

  app_settings = {
    "tftest" : "tftestvalue",
    "tftestkvp1" : "tftestkvpvalue1"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) python(data acceptance.TestData, pythonVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%s
resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "python"
  runtime_version                        = "%s"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, pythonVersion)
}

func (r FunctionAppFlexConsumptionResource) java(data acceptance.TestData, javaVersion string) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "java"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, javaVersion)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "java"
  runtime_version                        = "%s"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, javaVersion)
}

func (r FunctionAppFlexConsumptionResource) dotNet(data acceptance.TestData, dotNetVersion string) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "dotnet-isolated"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, dotNetVersion)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "dotnet-isolated"
  runtime_version                        = "%s"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, dotNetVersion)
}

func (r FunctionAppFlexConsumptionResource) powerShell(data acceptance.TestData, powerShellVersion string) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                = "powershell"
  runtime_version             = "%s"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, powerShellVersion)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "powershell"
  runtime_version                        = "%s"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, powerShellVersion)
}

func (r FunctionAppFlexConsumptionResource) runtimeNode(data acceptance.TestData, nodeVersion string) string {
	if !features.FivePointOh() {
		return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  storage_container_type            = "blobContainer"
  storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  storage_authentication_type       = "UserAssignedIdentity"
  storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                      = "node"
  runtime_version                   = "%s"
  maximum_instance_count            = 50
  instance_memory_in_mb             = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, nodeVersion)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "%s"
  maximum_instance_count                       = 50
  instance_memory_in_mb                        = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger, nodeVersion)
}

func (r FunctionAppFlexConsumptionResource) storageSystemAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "SystemAssignedIdentity"
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 50
  instance_memory_in_mb                  = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageUserAssignedIdentity1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s
resource "azurerm_user_assigned_identity" "test1" {
  name                = "acctest-uai1-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test1.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 50
  instance_memory_in_mb                        = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) storageUserAssignedIdentity2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 50
  instance_memory_in_mb                        = 2048

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) instanceMemoryUpdate(data acceptance.TestData, nodeVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                = "acctest-LFA-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "%s"
  maximum_instance_count                       = 50
  instance_memory_in_mb                        = 4096

  site_config {}
}
`, r.template(data), data.RandomInteger, nodeVersion)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = azurerm_storage_container.test.id
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 100
  instance_memory_in_mb                        = 2048
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = azurerm_storage_container.test.id
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 100
  instance_memory_in_mb                        = 2048
  always_ready {
    name           = "blob"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCountError(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = azurerm_storage_container.test.id
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 50
  instance_memory_in_mb                        = 2048
  always_ready {
    name           = "blob"
    instance_count = 20
  }
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 50
  }

  site_config {}
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) alwaysReadyInstanceCountUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctest-uai2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type            = "blobContainer"
  deployment_storage_container_endpoint        = azurerm_storage_container.test.id
  deployment_storage_authentication_type       = "UserAssignedIdentity"
  deployment_storage_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  runtime_name                                 = "node"
  runtime_version                              = "20"
  maximum_instance_count                       = 100
  instance_memory_in_mb                        = 2048
  always_ready {
    name           = "function:myHelloWorldFunction"
    instance_count = 20
  }
  always_ready {
    name           = "blob"
    instance_count = 20
  }

  site_config {
    application_insights_key               = azurerm_application_insights.test.instrumentation_key
    application_insights_connection_string = azurerm_application_insights.test.connection_string
  }
}
`, r.template(data), data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) vNetIntegration_subnetWithVnetProperties(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "acctest-subnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }

  lifecycle {
    ignore_changes = [
      delegation[0].service_delegation[0].actions
    ]
  }
}

resource "azurerm_function_app_flex_consumption" "test" {
  name                       = "acctest-LFA-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  virtual_network_subnet_id  = azurerm_subnet.test1.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  deployment_storage_container_type      = "blobContainer"
  deployment_storage_container_endpoint  = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}"
  deployment_storage_authentication_type = "StorageAccountConnectionString"
  deployment_storage_access_key          = azurerm_storage_account.test.primary_access_key
  runtime_name                           = "node"
  runtime_version                        = "20"
  maximum_instance_count                 = 100
  instance_memory_in_mb                  = 2048

  site_config {
    vnet_route_all_enabled = true
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r FunctionAppFlexConsumptionResource) templateKeyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acct-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

`, r.template(data), data.RandomInteger)
}

func (FunctionAppFlexConsumptionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-LFA-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "test" {
  name                = "acctestappinsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "acctestblobforfc"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_container" "test1" {
  name                  = "acctestblobforfc1"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_account" "test1" {
  name                     = "acctestsa1%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test1-1" {
  name                  = "acctestblobforfc11"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2021-04-01"
  expiry = "2024-03-30"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "FC1"
}

`, data.RandomInteger, "eastus2", data.RandomInteger, data.RandomInteger, data.RandomString, data.RandomString, data.RandomInteger) // location needs to be hardcoded for the moment because flex isn't available in all regions yet and appservice already has location overrides in TC
}
