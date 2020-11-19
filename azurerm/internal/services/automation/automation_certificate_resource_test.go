package automation_test

import (
	`context`
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure`
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationCertificateResource struct {
}

var (
	testCertThumbprintRaw, _ = ioutil.ReadFile(filepath.Join("testdata", "automation_certificate_test.thumb"))
	testCertRaw, _           = ioutil.ReadFile(filepath.Join("testdata", "automation_certificate_test.pfx"))
)

var testCertBase64 = base64.StdEncoding.EncodeToString(testCertRaw)

func TestAccAzureRMAutomationCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_certificate", "test")
	r := AutomationCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("base64"),
	})
}

func TestAccAzureRMAutomationCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_certificate", "test")
	r := AutomationCertificateResource{}

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

func TestAccAzureRMAutomationCertificate_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_certificate", "test")
	r := AutomationCertificateResource{}
	testCertThumbprint := strings.TrimSpace(string(testCertThumbprintRaw))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("base64").HasValue(testCertBase64),
				check.That(data.ResourceName).Key("thumbprint").HasValue(testCertThumbprint),
			),
		},
		data.ImportStep("base64"),
	})
}

func TestAccAzureRMAutomationCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_certificate", "test")
	r := AutomationCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue(""),
			),
		},
		data.ImportStep("base64"),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("description").HasValue("This is a test certificate for terraform acceptance test"),
			),
		},
		data.ImportStep("base64"),
	})
}

func (t AutomationCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {

	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	accountName := id.Path["automationAccounts"]
	name := id.Path["certificates"]

	resp, err := clients.Automation.CertificateClient.Get(ctx, id.ResourceGroup, accountName, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Certificate %q (resource group: %q): %+v", name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.CertificateProperties != nil), nil
}

func (AutomationCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, testCertBase64)
}

func (AutomationCertificateResource) requiresImport(data acceptance.TestData) string {
	template := AutomationCertificateResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_certificate" "import" {
  name                    = azurerm_automation_certificate.test.name
  resource_group_name     = azurerm_automation_certificate.test.resource_group_name
  automation_account_name = azurerm_automation_certificate.test.automation_account_name
  base64                  = azurerm_automation_certificate.test.base64
}
`, template)
}

func (AutomationCertificateResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
  description             = "This is a test certificate for terraform acceptance test"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, testCertBase64)
}

func (AutomationCertificateResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_certificate" "test" {
  name                    = "acctest-%d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  base64                  = "%s"
  description             = "This is a test certificate for terraform acceptance test"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, testCertBase64)
}
