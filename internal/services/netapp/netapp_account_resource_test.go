// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppAccountResource struct{}

func TestAccNetAppAccount(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests since
	// Azure allows only one active directory can be joined to a single subscription at a time for NetApp Account.
	// The CI system runs all tests in parallel, so the tests need to be changed to run one at a time.
	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":          testAccNetAppAccount_basic,
			"requiresImport": testAccNetAppAccount_requiresImport,
			"complete":       testAccNetAppAccount_complete,
			"update":         testAccNetAppAccount_update,
		},
	}

	for group, m := range testCases {
		for name, tc := range m {
			t.Run(group, func(t *testing.T) {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			})
		}
	}
}

func testAccNetAppAccount_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccNetAppAccount_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_netapp_account"),
		},
	})
}

func testAccNetAppAccount_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
			),
		},
		data.ImportStep("active_directory"),
	})
}

func testAccNetAppAccount_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_directory.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		{
			Config: r.completeConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("active_directory.#").HasValue("1"),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
			),
		},
		data.ImportStep("active_directory"),
	})
}

func TestAccNetAppAccount_systemAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccount_userAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.userAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppAccount_updateManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_account", "test")
	r := NetAppAccountResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.systemAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
			),
		},
		data.ImportStep(),
		{
			Config: r.userAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppAccountResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := netappaccounts.ParseNetAppAccountID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.AccountClient.AccountsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Account (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r NetAppAccountResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z",
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppAccountResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_account" "import" {
  name                = azurerm_netapp_account.test.name
  location            = azurerm_netapp_account.test.location
  resource_group_name = azurerm_netapp_account.test.resource_group_name

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z",
  }
}
`, r.basicConfig(data))
}

func (r NetAppAccountResource) completeConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns_servers         = ["1.2.3.4"]
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z",
    "FoO"           = "BaR"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppAccountResource) systemAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z",
    "FoO"           = "BaR"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r NetAppAccountResource) userAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}
`, r.template(data), data.RandomInteger)
}

func (NetAppAccountResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z",
    "SkipNRMSNSG"   = "true"
  }
}


`, data.RandomInteger, data.Locations.Primary)
}
