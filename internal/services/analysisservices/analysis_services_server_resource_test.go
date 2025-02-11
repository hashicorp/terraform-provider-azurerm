// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package analysisservices_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AnalysisServicesServerResource struct{}

func TestAccAnalysisServicesServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

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

func TestAccAnalysisServicesServer_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.querypoolConnectionMode(data, "All"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("querypool_connection_mode").HasValue("All"),
			),
		},
		data.ImportStep(),
		{
			Config: r.querypoolConnectionMode(data, "ReadOnly"),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.firewallSettings1(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ipv4_firewall_rule.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.firewallSettings2(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ipv4_firewall_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.firewallSettings3(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
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
		t.Skipf("Acceptance test skipped unless env '%s' and '%s' set", ArmAccAdminEmail1, ArmAccAdminEmail2)
		return
	}

	email1 := os.Getenv(ArmAccAdminEmail1)
	email2 := os.Getenv(ArmAccAdminEmail2)
	preAdminUsers := []string{email1}
	postAdminUsers := []string{email1, email2}

	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.adminUsers(data, preAdminUsers),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.adminUsers(data, postAdminUsers),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAzureRMAnalysisServicesServer_serverFullName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_analysis_services_server", "test")
	r := AnalysisServicesServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serverFullName(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupBlobContainerUri(data),
			Check: acceptance.ComposeTestCheckFunc(
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
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				data.CheckWithClient(r.suspend),
				data.CheckWithClient(r.checkState(servers.StatePaused)),
			),
		},
		data.ImportStep(),
		{
			Config: r.scale(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("B2"),
				data.CheckWithClient(r.checkState(servers.StatePaused)),
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

func (t AnalysisServicesServerResource) withTags(data acceptance.TestData) string {
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
  name     = "acctestRG-analysis-%d"
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
  name     = "acctestRG-analysis-%d"
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
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                     = "acctestass%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "B1"
  power_bi_service_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enablePowerBIService)
}

func (t AnalysisServicesServerResource) firewallSettings2(data acceptance.TestData, enablePowerBIService bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                     = "acctestass%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "B1"
  power_bi_service_enabled = %t

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
  name     = "acctestRG-analysis-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                     = "acctestass%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku                      = "B1"
  power_bi_service_enabled = %t

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
  name     = "acctestRG-analysis-%d"
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
  sku                 = "B2"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t AnalysisServicesServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servers.ParseServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AnalysisServices.Servers.GetDetails(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (t AnalysisServicesServerResource) suspend(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	client := clients.AnalysisServices.Servers

	id, err := servers.ParseServerID(state.ID)
	if err != nil {
		return err
	}

	timeout, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	if err := client.SuspendThenPoll(timeout, *id); err != nil {
		return fmt.Errorf("suspending %s: %+v", *id, err)
	}

	return nil
}

func (t AnalysisServicesServerResource) checkState(expectedState servers.State) acceptance.ClientCheckFunc {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.AnalysisServices.Servers

		id, err := servers.ParseServerID(state.ID)
		if err != nil {
			return err
		}

		timeout, cancel := context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()

		resp, err := client.GetDetails(timeout, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s to check the state: %+v", *id, err)
		}

		actualState := ""
		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
			actualState = string(*resp.Model.Properties.State)
		}

		if actualState != string(expectedState) {
			return fmt.Errorf("Unexpected state. Expected %s but is %s", string(expectedState), actualState)
		}

		return nil
	}
}
