package resourceconnector_test

import (
	"context"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceConnectorApplianceResource struct{}

func TestAccResourceConnectorAppliance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	r := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceConnectorAppliance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	r := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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

func TestAccResourceConnectorAppliance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	r := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, r.generatePublicKey()),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceConnectorAppliance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	r := ResourceConnectorApplianceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
		},
		{
			Config:      r.requiresImport(data, r.generatePublicKey()),
			ExpectError: acceptance.RequiresImportError("azurerm_resource_connector_appliance"),
		},
	})
}

func (r ResourceConnectorApplianceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appliances.ParseApplianceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ResourceConnector.AppliancesClient
	resp, err := client.Get(ctx, *id)

	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ResourceConnectorApplianceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestrcapplicance-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ResourceConnectorApplianceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestrcapplicance-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  identity {
    type = "SystemAssigned"
  }
  tags = {
    "hello" = "world"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r ResourceConnectorApplianceResource) complete(data acceptance.TestData, publicKey string) string {
	return fmt.Sprintf(`

%s

resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestrcapplicance-%[2]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  public_key              = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  tags = {
    "hello" = "world"
  }
}
`, r.template(data), data.RandomInteger, publicKey)
}

func (r ResourceConnectorApplianceResource) requiresImport(data acceptance.TestData, publicKey string) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_resource_connector_appliance" "import" {
  name                    = azurerm_resource_connector_appliance.test.name
  location                = azurerm_resource_connector_appliance.test.location
  resource_group_name     = azurerm_resource_connector_appliance.test.resource_group_name
  distro                  = azurerm_resource_connector_appliance.test.distro
  infrastructure_provider = azurerm_resource_connector_appliance.test.infrastructure_provider
  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data), data.RandomInteger)
}

func (r ResourceConnectorApplianceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-appliances-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceConnectorApplianceResource) generatePublicKey() string {
	privateKey, err := rsa.GenerateKey(cryptoRand.Reader, 4096)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))
}
