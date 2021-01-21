package resource_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azureResource "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ManagementLockResource struct {
}

func TestAccManagementLock_resourceGroupReadOnlyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupReadOnlyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupReadOnlyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccManagementLock_resourceGroupReadOnlyComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupReadOnlyComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_resourceGroupCanNotDeleteBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupCanNotDeleteBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_resourceGroupCanNotDeleteComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.resourceGroupCanNotDeleteComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_publicIPReadOnlyBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.publicIPReadOnlyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_publicIPCanNotDeleteBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.publicIPCanNotDeleteBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_subscriptionReadOnlyBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionReadOnlyBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementLock_subscriptionCanNotDeleteBasic(t *testing.T) {
	_, exists := os.LookupEnv("TF_ACC_SUBSCRIPTION_PARALLEL_LOCK")
	if !exists {
		t.Skip("`TF_ACC_SUBSCRIPTION_PARALLEL_LOCK` isn't specified - skipping since this test can't be run in Parallel")
	}

	data := acceptance.BuildTestData(t, "azurerm_management_lock", "test")
	r := ManagementLockResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionCanNotDeleteBasic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t ManagementLockResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azureResource.ParseAzureRMLockId(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.LocksClient.GetByScope(ctx, id.Scope, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Management Lock (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (ManagementLockResource) resourceGroupReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ManagementLockResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_lock" "import" {
  name       = azurerm_management_lock.test.name
  scope      = azurerm_management_lock.test.scope
  lock_level = azurerm_management_lock.test.lock_level
}
`, r.resourceGroupReadOnlyBasic(data))
}

func (ManagementLockResource) resourceGroupReadOnlyComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "ReadOnly"
  notes      = "Hello, World!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagementLockResource) resourceGroupCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagementLockResource) resourceGroupCanNotDeleteComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_resource_group.test.id
  lock_level = "CanNotDelete"
  notes      = "Hello, World!"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagementLockResource) publicIPReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_public_ip.test.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagementLockResource) publicIPCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = azurerm_public_ip.test.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ManagementLockResource) subscriptionReadOnlyBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = data.azurerm_subscription.current.id
  lock_level = "ReadOnly"
}
`, data.RandomInteger)
}

func (ManagementLockResource) subscriptionCanNotDeleteBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {
}

resource "azurerm_management_lock" "test" {
  name       = "acctestlock-%d"
  scope      = data.azurerm_subscription.current.id
  lock_level = "CanNotDelete"
}
`, data.RandomInteger)
}
