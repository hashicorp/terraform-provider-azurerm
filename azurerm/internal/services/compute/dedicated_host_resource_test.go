package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DedicatedHostResource struct {
}

func TestAccDedicatedHost_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_basicNewSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicNewSku(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_autoReplaceOnFailure(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			// Enabled
			Config: r.autoReplaceOnFailure(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Disabled
			Config: r.autoReplaceOnFailure(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// Enabled
			Config: r.autoReplaceOnFailure(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_licenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.licenceType(data, "None"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenceType(data, "Windows_Server_Hybrid"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenceType(data, "Windows_Server_Perpetual"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.licenceType(data, "None"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDedicatedHost_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dedicated_host", "test")
	r := DedicatedHostResource{}

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

func (t DedicatedHostResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DedicatedHostID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DedicatedHostsClient.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Dedicated Host %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r DedicatedHostResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}
`, r.template(data), data.RandomInteger)
}

func (r DedicatedHostResource) basicNewSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DCSv2-Type1"
  platform_fault_domain   = 1
}
`, r.template(data), data.RandomInteger)
}

func (r DedicatedHostResource) autoReplaceOnFailure(data acceptance.TestData, replace bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  auto_replace_on_failure = %t
}
`, r.template(data), data.RandomInteger, replace)
}

func (r DedicatedHostResource) licenceType(data acceptance.TestData, licenseType string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  license_type            = %q
}
`, r.template(data), data.RandomInteger, licenseType)
}

func (r DedicatedHostResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dedicated_host" "test" {
  name                    = "acctest-DH-%d"
  location                = azurerm_resource_group.test.location
  dedicated_host_group_id = azurerm_dedicated_host_group.test.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
  license_type            = "Windows_Server_Hybrid"
  auto_replace_on_failure = false
}
`, r.template(data), data.RandomInteger)
}

func (r DedicatedHostResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host" "import" {
  name                    = azurerm_dedicated_host.test.name
  location                = azurerm_dedicated_host.test.location
  dedicated_host_group_id = azurerm_dedicated_host.test.dedicated_host_group_id
  sku_name                = azurerm_dedicated_host.test.sku_name
  platform_fault_domain   = azurerm_dedicated_host.test.platform_fault_domain
}
`, r.basic(data))
}

func (DedicatedHostResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-compute-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  name                        = "acctest-DHG-%d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  platform_fault_domain_count = 2
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
