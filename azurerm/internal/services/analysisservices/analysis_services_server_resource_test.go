package analysisservices_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/analysisservices/mgmt/2017-08-01/analysisservices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/analysisservices/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AnalysisServicesServerResource struct {
}

func TestAccAnalysisServicesServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

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

func TestAccAnalysisServicesServer_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAnalysisServicesServer_querypoolConnectionMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.querypoolConnectionMode(data, "All"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("querypool_connection_mode").HasValue("All"),
			),
		},
		data.ImportStep(),
		{
			Config: r.querypoolConnectionMode(data, "ReadOnly"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("querypool_connection_mode").HasValue("ReadOnly"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAnalysisServicesServer_firewallSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.firewallSettings1(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_power_bi_service").HasValue("true"),
				check.That(data.ResourceName).Key("ipv4_firewall_rule.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.firewallSettings2(data, false),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_power_bi_service").HasValue("false"),
				check.That(data.ResourceName).Key("ipv4_firewall_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.firewallSettings3(data, true),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enable_power_bi_service").HasValue("true"),
				check.That(data.ResourceName).Key("ipv4_firewall_rule.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

// ARM_ACC_EMAIL1 and ARM_ACC_EMAIL2 must be set and existing emails in the tenant's AD to work properly
func TestAccAzureRMAnalysisServicesServer_adminUsers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")

	const ArmAccAdminEmail1 = "ARM_ACCTEST_ADMIN_EMAIL1"
	const ArmAccAdminEmail2 = "ARM_ACCTEST_ADMIN_EMAIL2"

	if os.Getenv(ArmAccAdminEmail1) == "" || os.Getenv(ArmAccAdminEmail2) == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env '%s' and '%s' set", ArmAccAdminEmail1, ArmAccAdminEmail2))
		return
	}

	email1 := os.Getenv(ArmAccAdminEmail1)
	email2 := os.Getenv(ArmAccAdminEmail2)
	preAdminUsers := []string{email1}
	postAdminUsers := []string{email1, email2}

	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.adminUsers(data, preAdminUsers),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.adminUsers(data, postAdminUsers),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAnalysisServicesServer_serverFullName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.serverFullName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("server_full_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAnalysisServicesServer_backupBlobContainerUri(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.backupBlobContainerUri(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup_blob_container_uri").Exists(),
			),
		},
		data.ImportStep("backup_blob_container_uri"),
	})
}

func TestAccAzureRMAnalysisServicesServer_suspended(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.suspend),
				data.CheckWithClient(r.checkState(analysisservices.StatePaused)),
			),
		},
		data.ImportStep(),
		{
			Config: r.scale(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("S1"),
				data.CheckWithClient(r.checkState(analysisservices.StatePaused)),
			),
		},
		data.ImportStep(),
	})
}

func (t AnalysisServicesServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) querypoolConnectionMode(data acceptance.TestData, connectionMode string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                      = "acctestass%d"
  location                  = azurerm_resource_group.test.location
  resource_group_name       = azurerm_resource_group.test.name
  sku                       = "B1"
  querypool_connection_mode = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, connectionMode)
}

func (t AnalysisServicesServerResource) firewallSettings1(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku                     = "B1"
  enable_power_bi_service = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func (t AnalysisServicesServerResource) firewallSettings2(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku                     = "B1"
  enable_power_bi_service = %t

  ipv4_firewall_rule {
    name        = "test1"
    range_start = "92.123.234.11"
    range_end   = "92.123.234.12"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func (t AnalysisServicesServerResource) firewallSettings3(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                    = "acctestass%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  sku                     = "B1"
  enable_power_bi_service = %t

  ipv4_firewall_rule {
    name        = "test1"
    range_start = "92.123.234.11"
    range_end   = "92.123.234.13"
  }

  ipv4_firewall_rule {
    name        = "test2"
    range_start = "226.202.187.57"
    range_end   = "226.208.192.47"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func (t AnalysisServicesServerResource) adminUsers(data acceptance.TestData, adminUsers []string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"
  admin_users         = ["%s"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, strings.Join(adminUsers, "\", \""))
}

func (t AnalysisServicesServerResource) serverFullName(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) backupBlobContainerUri(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestass%s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "assbackup"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

data "azurerm_storage_account_blob_container_sas" "test" {
  connection_string = azurerm_storage_account.test.primary_connection_string
  container_name    = azurerm_storage_container.test.name
  https_only        = true

  start  = "2018-06-01"
  expiry = "2048-06-01"

  permissions {
    read   = true
    add    = true
    create = true
    write  = true
    delete = true
    list   = true
  }
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "B1"

  backup_blob_container_uri = "${azurerm_storage_account.test.primary_blob_endpoint}${azurerm_storage_container.test.name}${data.azurerm_storage_account_blob_container_sas.test.sas}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func (t AnalysisServicesServerResource) scale(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "S1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AnalysisServices.ServerClient.GetDetails(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Analysis Services Server %q (resource group: %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ServerProperties != nil), nil
}

func (t AnalysisServicesServerResource) suspend(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
	client := clients.AnalysisServices.ServerClient

	id, err := parse.ServerID(state.ID)
	if err != nil {
		return err
	}

	suspendFuture, err := client.Suspend(ctx, id.ResourceGroup, id.ServerName)
	if err != nil {
		return fmt.Errorf("suspending Analysis Services Server %q (resource group: %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	err = suspendFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Wait for Suspend on Analysis Services Server %q (resource group: %q): %+v", id.ServerName, id.ResourceGroup, err)
	}

	return nil
}

func (t AnalysisServicesServerResource) checkState(serverState analysisservices.State) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) error {
		client := clients.AnalysisServices.ServerClient

		id, err := parse.ServerID(state.ID)
		if err != nil {
			return err
		}

		resp, err := client.GetDetails(ctx, id.ResourceGroup, id.ServerName)
		if err != nil {
			return fmt.Errorf("Bad: Get on analysisServicesServerClient: %+v", err)
		}

		if resp.State != serverState {
			return fmt.Errorf("Unexpected state. Expected %s but is %s", state, resp.State)
		}

		return nil
	}
}
