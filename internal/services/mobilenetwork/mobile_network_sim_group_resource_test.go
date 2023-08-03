// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/simgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MobileNetworkSimGroupResource struct{}

func TestAccMobileNetworkSimGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	r := MobileNetworkSimGroupResource{}
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

func TestAccMobileNetworkSimGroup_withEncryptionKeyUrl(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	r := MobileNetworkSimGroupResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withEncryptionKeyUrl(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMobileNetworkSimGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	r := MobileNetworkSimGroupResource{}
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

func TestAccMobileNetworkSimGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	r := MobileNetworkSimGroupResource{}
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

func TestAccMobileNetworkSimGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mobile_network_sim_group", "test")
	r := MobileNetworkSimGroupResource{}
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

func (r MobileNetworkSimGroupResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := simgroup.ParseSimGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.MobileNetwork.SIMGroupClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r MobileNetworkSimGroupResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_mobile_network_sim_group" "test" {
  name              = "acctest-mnsg-%d"
  location          = azurerm_resource_group.test.location
  mobile_network_id = azurerm_mobile_network.test.id
}
`, template, data.RandomInteger)
}

func (r MobileNetworkSimGroupResource) withEncryptionKeyUrl(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-mn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}


resource "azurerm_key_vault" "test" {
  name                = "acct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = data.azurerm_client_config.test.object_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "enckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_mobile_network_sim_group" "test" {
  name               = "acctest-mnsg-%[2]d"
  location           = "%[3]s"
  encryption_key_url = azurerm_key_vault_key.test.versionless_id
  mobile_network_id  = azurerm_mobile_network.test.id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r MobileNetworkSimGroupResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_mobile_network_sim_group" "import" {
  name              = azurerm_mobile_network_sim_group.test.name
  location          = azurerm_mobile_network_sim_group.test.location
  mobile_network_id = azurerm_mobile_network_sim_group.test.mobile_network_id

}
`, template)
}

func (r MobileNetworkSimGroupResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "test" {}

resource "azurerm_key_vault" "test" {
  name                = "acct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = data.azurerm_client_config.test.object_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
  }
}


resource "azurerm_key_vault_key" "test" {
  name         = "enckey%[2]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}


resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-mn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_mobile_network_sim_group" "test" {
  name               = "acctest-mnsg-%[2]d"
  location           = azurerm_mobile_network.test.location
  mobile_network_id  = azurerm_mobile_network.test.id
  encryption_key_url = azurerm_key_vault_key.test.versionless_id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r MobileNetworkSimGroupResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest-mn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_key_vault" "test" {
  name                = "acct-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.test.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = data.azurerm_client_config.test.object_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id          = data.azurerm_client_config.test.tenant_id
    object_id          = azurerm_user_assigned_identity.test.principal_id
    secret_permissions = ["Delete", "Get", "Set"]
    key_permissions    = ["Create", "Delete", "Get", "Import", "Purge", "UnwrapKey", "WrapKey", "GetRotationPolicy"]
  }
}


resource "azurerm_key_vault_key" "test" {
  name         = "enckey%[2]d"
  key_vault_id = "${azurerm_key_vault.test.id}"
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}


resource "azurerm_mobile_network_sim_group" "test" {
  name               = "acctest-mnsg-%[2]d"
  location           = azurerm_mobile_network.test.location
  mobile_network_id  = azurerm_mobile_network.test.id
  encryption_key_url = azurerm_key_vault_key.test.versionless_id

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  tags = {
    key = "updated"
  }
}
`, template, data.RandomInteger)
}

func (r MobileNetworkSimGroupResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-mn-%[1]d"
  location = %[2]q
}

resource "azurerm_mobile_network" "test" {
  name                = "acctest-mn-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  mobile_country_code = "001"
  mobile_network_code = "01"
}
`, data.RandomInteger, data.Locations.Primary)
}
