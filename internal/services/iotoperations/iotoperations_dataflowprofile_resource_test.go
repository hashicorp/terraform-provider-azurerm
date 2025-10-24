package iotoperations_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/dataflowprofile"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsDataflowProfileResource is a test harness for azurerm_iotoperations_dataflow_profile acceptance tests.
type IotOperationsDataflowProfileResource struct{}

// NOTE: These tests currently use a placeholder instance name "REPLACE_WITH_INSTANCE_NAME".
// To make them pass, add or reference an actual azurerm_iotoperations_instance in the template.

func TestAccIotOperationsDataflowProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

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

func TestAccIotOperationsDataflowProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotoperations_dataflow_profile"),
		},
	})
}

func TestAccIotOperationsDataflowProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

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

func TestAccIotOperationsDataflowProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotoperations_dataflow_profile", "test")
	r := IotOperationsDataflowProfileResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
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

// Exists implements the acceptance existence check using the generated SDK ID parser.
func (IotOperationsDataflowProfileResource) Exists(ctx context.Context, c *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dataflowprofile.ParseDataflowProfileID(state.ID)
	if err != nil {
		return nil, fmt.Errorf("parsing ID %q: %w", state.ID, err)
	}

	resp, err := c.IoTOperations.DataflowProfileClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id.ID(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

// template builds the minimal provider + resource_group.
// TODO: Create or reference a real azurerm_iotoperations_instance and feed its name into the dataflow profile.
func (IotOperationsDataflowProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
}
`, data.RandomInteger)
}

func (r IotOperationsDataflowProfileResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# TODO: create or reference an IoT Operations instance here (azurerm_iotoperations_instance).

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-dfp-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "REPLACE_WITH_INSTANCE_NAME"

  properties {
    instance_count = 1
  }
}
`, r.template(data), data.RandomString)
}

func (r IotOperationsDataflowProfileResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotoperations_dataflow_profile" "import" {
  name                = azurerm_iotoperations_dataflow_profile.test.name
  resource_group_name = azurerm_iotoperations_dataflow_profile.test.resource_group_name
  instance_name       = azurerm_iotoperations_dataflow_profile.test.instance_name

  properties {
    instance_count = 1
  }
}
`, r.basic(data))
}

func (r IotOperationsDataflowProfileResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

# TODO: replace the placeholder instance name with an actual one.

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-dfp-%s"
  resource_group_name = azurerm_resource_group.test.name
  instance_name       = "REPLACE_WITH_INSTANCE_NAME"
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }

  properties {
    instance_count = 3

    diagnostics {
      logs {
        level = "Info"
      }
      metrics {
        prometheus_port = 7581
      }
    }
  }
}
`, r.template(data), data.RandomString)
}
