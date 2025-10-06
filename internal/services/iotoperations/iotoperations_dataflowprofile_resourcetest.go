package iotoperations_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// IotOperationsDataflowProfileResource is a test harness for azurerm_iotoperations_dataflow_profile acceptance tests.
type IotOperationsDataflowProfileResource struct{}

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

func (IotOperationsDataflowProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    rg, instance, profile, err := parseDataflowProfileID(state.ID)
    if err != nil {
        return nil, err
    }

    resp, err := clients.IoTOperations.DataflowProfileClient.Get(ctx, rg, instance, profile, nil)
    if err != nil {
        return nil, fmt.Errorf("retrieving %s: %+v", state.ID, err)
    }

    return utils.Bool(resp.DataflowProfileResource != nil), nil
}

// template builds the minimal provider + resource_group; NOTE: you must create an IoT Operations instance
// resources in the template below or reference pre-existing ones. 
func (IotOperationsDataflowProfileResource) template(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-iotops-%d"
  location = "%s"
}
`, data.RandomInteger, data.Location())
}

func (r IotOperationsDataflowProfileResource) basic(data acceptance.TestData) string {
    // TODO: Replace placeholders with a real IoT Operations instance resource (azurerm_iotoperations_instance).
    return fmt.Sprintf(`
%s

# TODO: create or reference an IoT Operations instance here (azurerm_iotoperations_instance).

resource "azurerm_iotoperations_dataflow_profile" "test" {
  name                = "acctest-dfp-%s"
  resource_group_name = azurerm_resource_group.test.name

  # Replace with the actual instance name created in this template:
  instance_name = "REPLACE_WITH_INSTANCE_NAME"

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

# See TODO in basic template for instance creation.

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

// parseDataflowProfileID extracts resource group, instance name and profile name from a full ARM resource id.
// Expected pattern:
// /subscriptions/.../resourceGroups/{rg}/providers/Microsoft.IoTOperations/instances/{instance}/dataflowProfiles/{profile}
func parseDataflowProfileID(id string) (rg, instance, profile string, err error) {
    if id == "" {
        return "", "", "", fmt.Errorf("empty id")
    }
    parts := strings.Split(id, "/")
    // Normalize: drop leading empty elements if id starts with '/'
    start := 0
    if parts[0] == "" {
        start = 1
    }
    for i := start; i < len(parts)-1; i++ {
        switch strings.ToLower(parts[i]) {
        case "resourcegroups":
            rg = parts[i+1]
        case "instances":
            instance = parts[i+1]
        case "dataflowprofiles":
            profile = parts[i+1]
        }
    }
    if rg == "" || instance == "" || profile == "" {
        return "", "", "", fmt.Errorf("failed to parse id: %s", id)
    }
    return rg, instance, profile, nil
}