package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineGroupResource struct{}

func TestAccMsSqlVirtualMachineGroupResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_group", "test")
	r := MsSqlVirtualMachineGroupResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMsSqlVirtualMachineGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_group", "test")
	r := MsSqlVirtualMachineGroupResource{}
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

func TestAccMsSqlVirtualMachineGroupResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_group", "test")
	r := MsSqlVirtualMachineGroupResource{}

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

func (MsSqlVirtualMachineGroupResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.VirtualMachineGroupsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}
		return nil, fmt.Errorf("reading %s: %v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MsSqlVirtualMachineGroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_virtual_machine_group" "test" {
  name                = "acctestag%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sql_image_offer = "SQL2017-WS2016"
  sql_image_sku   = "Developer"

  wsfc_domain_profile {
    fqdn                = "testdomain.com"
    cluster_subnet_type = "SingleSubnet"
  }
}




`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func (r MsSqlVirtualMachineGroupResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_mssql_virtual_machine_group" "import" {
  name                = azurerm_mssql_virtual_machine_group.test.name
  resource_group_name = azurerm_mssql_virtual_machine_group.test.resource_group_name
  location            = azurerm_mssql_virtual_machine_group.test.location

  sql_image_offer = "SQL2017-WS2016"
  sql_image_sku   = "Developer"

  wsfc_domain_profile {
    fqdn                = "testdomain.com"
    cluster_subnet_type = "SingleSubnet"
  }
}
`, r.basic(data))
}

func (r MsSqlVirtualMachineGroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "unlikely23exst2acct%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_virtual_machine_group" "test" {
  name                = "acctestag%[4]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sql_image_offer     = "SQL2017-WS2016"
  sql_image_sku       = "Developer"


  wsfc_domain_profile {
    fqdn                           = "testdomain.com"
    cluster_bootstrap_account_name = "bootstrapacc%[5]s"
    cluster_operator_account_name  = "opacc%[5]s"
    sql_service_account_name       = "sqlsrvacc%[5]s"
    storage_account_url            = azurerm_storage_account.test.primary_blob_endpoint
    cluster_subnet_type            = "SingleSubnet"
  }

  tags = {
    test = "testing"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomString)
}
