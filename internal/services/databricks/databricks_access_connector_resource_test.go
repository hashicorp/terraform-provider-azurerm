package databricks_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2022-10-01-preview/accessconnector"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DatabricksAccessConnectorResource struct{}

func TestAccDatabricksAccessConnector_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}

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

func TestAccDatabricksAccessConnector_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}

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

func TestAccDatabricksAccessConnector_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDatabricksAccessConnector_identityComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.identityUserAssigned(data),
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

func TestAccDatabricksAccessConnector_identityMissingIdentityIDsError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.identityUserAssignedMissingIdentityIDs(data),
			Check:       acceptance.ComposeTestCheckFunc(),
			ExpectError: regexp.MustCompile(`must be specified when`),
		},
	})
}

func TestAccDatabricksAccessConnector_identityUserAssignedTooManyError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_databricks_access_connector", "test")
	r := DatabricksAccessConnectorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.identityUserAssignedTooManyError(data),
			Check:       acceptance.ComposeTestCheckFunc(),
			ExpectError: regexp.MustCompile(`Too many list items`),
		},
	})
}

func (DatabricksAccessConnectorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := accessconnector.ParseAccessConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DataBricks.AccessConnectorClient.Get(ctx, *id)

	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}

	if err != nil {
		return nil, fmt.Errorf("making Read request on Databricks %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DatabricksAccessConnectorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-databricks-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DatabricksAccessConnectorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBAC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r DatabricksAccessConnectorResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_databricks_access_connector" "import" {
  name                = azurerm_databricks_access_connector.test.name
  resource_group_name = azurerm_databricks_access_connector.test.resource_group_name
  location            = azurerm_databricks_access_connector.test.location
}
`, template)
}

func (r DatabricksAccessConnectorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBAC%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r DatabricksAccessConnectorResource) identityUserAssigned(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestDBUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBAC%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"

    identity_ids = [
      azurerm_user_assigned_identity.test.id,
    ]
  }
}
`, template, data.RandomInteger)
}

func (r DatabricksAccessConnectorResource) identityUserAssignedMissingIdentityIDs(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestDBUAI-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBAC%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"
  }
}
`, template, data.RandomInteger)
}

func (r DatabricksAccessConnectorResource) identityUserAssignedTooManyError(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_user_assigned_identity" "R2" {
  name                = "acctestUAIR2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_user_assigned_identity" "D2" {
  name                = "acctestUAID2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_databricks_access_connector" "test" {
  name                = "acctestDBAC-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.R2.id,
      azurerm_user_assigned_identity.D2.id,
    ]
  }
}
`, template, data.RandomInteger)
}
