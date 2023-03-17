package hybridcompute_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/machineextensions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HybridComputeMachineExtensionResource struct{}

func TestAccHybridComputeMachineExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("publisher").HasValue("Microsoft.Azure.Monitor"),
				check.That(data.ResourceName).Key("type").HasValue("AzureMonitorLinuxAgent"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHybridComputeMachineExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
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

func TestAccHybridComputeMachineExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("publisher").HasValue("Microsoft.Azure.Monitor"),
				check.That(data.ResourceName).Key("type").HasValue("AzureMonitorLinuxAgent"),
				check.That(data.ResourceName).Key("type_handler_version").MatchesRegex(regexp.MustCompile("^1[.]24.*$")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccHybridComputeMachineExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
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
	})
}

func (r HybridComputeMachineExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := machineextensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.HybridCompute.MachineExtensionsClient
	resp, err := client.Get(ctx, *id)
	exists := false
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return &exists, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exists = resp.Model != nil
	return &exists, nil
}

func (r HybridComputeMachineExtensionResource) template(data acceptance.TestData) string {
	d := HybridComputeMachineDataSource{}
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	randomUUID, _ := uuid.GenerateUUID()
	password := generateRandomPassword(10)
	return d.basic(data, clientSecret, randomUUID, password)
}

func (r HybridComputeMachineExtensionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                      = "acctest-hcme-%d"
  hybrid_compute_machine_id = data.azurerm_hybrid_compute_machine.test.id
  publisher                 = "Microsoft.Azure.Monitor"
  type                      = "AzureMonitorLinuxAgent"
  location                  = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r HybridComputeMachineExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "import" {
  name                      = azurerm_hybrid_compute_machine_extension.test.name
  hybrid_compute_machine_id = azurerm_hybrid_compute_machine_extension.test.id
  publisher                 = azurerm_hybrid_compute_machine_extension.test.publisher
  type                      = azurerm_hybrid_compute_machine_extension.test.type
  location                  = azurerm_hybrid_compute_machine_extension.test.location
}
`, config)
}

func (r HybridComputeMachineExtensionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                               = "acctest-hcme-%d"
  hybrid_compute_machine_id          = data.azurerm_hybrid_compute_machine.test.id
  location                           = "%s"
  auto_upgrade_minor_version_enabled = false
  automatic_upgrade_enabled          = false
  publisher                          = "Microsoft.Azure.Monitor"
  type                               = "AzureMonitorLinuxAgent"
  type_handler_version               = "1.24"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r HybridComputeMachineExtensionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                               = "acctest-hcme-%d"
  hybrid_compute_machine_id          = data.azurerm_hybrid_compute_machine.test.id
  location                           = "%s"
  auto_upgrade_minor_version_enabled = false
  automatic_upgrade_enabled          = true
  publisher                          = "Microsoft.Azure.Monitor"
  type                               = "AzureMonitorLinuxAgent"
  type_handler_version               = "1.25"
}
`, template, data.RandomInteger, data.Locations.Primary)
}
