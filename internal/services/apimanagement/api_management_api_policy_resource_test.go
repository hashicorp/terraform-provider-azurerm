// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apipolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementApiPolicyResource struct{}

func TestAccApiManagementAPIPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "test")
	r := ApiManagementApiPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"xml_link"},
		},
	})
}

func TestAccApiManagementAPIPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "test")
	r := ApiManagementApiPolicyResource{}

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

func TestAccApiManagementAPIPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "test")
	r := ApiManagementApiPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
		{
			Config: r.customPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.importStep(data, r.customPolicyXml()),
	})
}

func TestAccApiManagementAPIPolicy_customPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "test")
	r := ApiManagementApiPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		r.importStep(data, r.customPolicyXml()),
	})
}

func TestAccApiManagementAPIPolicy_rawXmlUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_policy", "test")
	r := ApiManagementApiPolicyResource{}

	policy, err := os.ReadFile("testdata/api_management_api_policy_raw.xml")
	if err != nil {
		t.Fatal(err)
	}

	updatedPolicy, err := os.ReadFile("testdata/api_management_api_policy_raw_updated.xml")
	if err != nil {
		t.Fatal(err)
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.rawXml(data, "api_management_api_policy_raw.xml"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("xml_content").HasValue(string(policy)),
			),
		},
		r.importStep(data, string(policy)),
		{
			Config: r.rawXml(data, "api_management_api_policy_raw_updated.xml"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("xml_content").HasValue(string(updatedPolicy)),
			),
		},
		r.importStep(data, string(updatedPolicy)),
		{
			Config: r.rawXml(data, "api_management_api_policy_raw.xml"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("xml_content").HasValue(string(policy)),
			),
		},
		r.importStep(data, string(policy)),
	})
}

func (ApiManagementApiPolicyResource) importStep(data acceptance.TestData, policyContent string) acceptance.TestStep {
	step := data.ImportStep("xml_content", "xml_link")
	step.ImportStateCheck = func(states []*pluginsdk.InstanceState) error {
		if len(states) != 1 {
			return fmt.Errorf("expected one imported policy, got %d", len(states))
		}

		importedContent := states[0].Attributes["xml_content"]
		if !apimanagement.XmlWithDotNetInterpolationsDiffSuppress("xml_content", policyContent, importedContent, nil) {
			return fmt.Errorf("imported xml_content differs from the expected policy: expected %q, got %q", policyContent, importedContent)
		}

		return nil
	}
	return step
}

func (ApiManagementApiPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apipolicy.ParseApiID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiPoliciesClient.Get(ctx, *id, apipolicy.GetOperationOptions{Format: pointer.To(apipolicy.PolicyExportFormatXml)})
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (ApiManagementApiPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}

resource "azurerm_api_management_api_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  xml_link            = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/refs/heads/main/internal/services/apimanagement/testdata/api_management_policy_test.xml"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementApiPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_policy" "import" {
  api_name            = azurerm_api_management_api_policy.test.api_name
  api_management_name = azurerm_api_management_api_policy.test.api_management_name
  resource_group_name = azurerm_api_management_api_policy.test.resource_group_name
  xml_link            = azurerm_api_management_api_policy.test.xml_link
}
`, r.basic(data))
}

func (r ApiManagementApiPolicyResource) customPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Consumption_0"
}

resource "azurerm_api_management_api" "test" {
  name                = "acctestapi-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "api1"
  path                = "api1"
  protocols           = ["https"]
  revision            = "1"
}

resource "azurerm_api_management_api_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name

  xml_content = <<XML
%s
XML

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, r.customPolicyXml())
}

func (ApiManagementApiPolicyResource) customPolicyXml() string {
	return `<policies>
  <inbound>
    <set-variable name="abc" value="@(context.Request.Headers.GetValueOrDefault("X-Header-Name", ""))" />
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>`
}

func (ApiManagementApiPolicyResource) rawXml(data acceptance.TestData, policyFile string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_policy" "test" {
  api_name            = azurerm_api_management_api.test.name
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  xml_content         = file("testdata/%s")
}
`, ApiManagementApiResource{}.basic(data), policyFile)
}
