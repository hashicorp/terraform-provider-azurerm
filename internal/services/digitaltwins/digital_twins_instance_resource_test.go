// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DigitalTwinsInstanceResource struct{}

func TestAccDigitalTwinsInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDigitalTwinsInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").DoesNotExist(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").DoesNotExist(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithIdentitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsUUID(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsUUID(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_identityUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithIdentityUserAssigned(data),
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
		{
			Config: r.basicWithIdentityUserAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (DigitalTwinsInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := digitaltwinsinstance.ParseDigitalTwinsInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DigitalTwins.InstanceClient.DigitalTwinsGet(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (DigitalTwinsInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dtwin-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DigitalTwinsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "import" {
  name                = azurerm_digital_twins_instance.test.name
  resource_group_name = azurerm_digital_twins_instance.test.resource_group_name
  location            = azurerm_digital_twins_instance.test.location
}
`, r.basic(data))
}

func (r DigitalTwinsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Stage"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) basicWithIdentitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) basicWithIdentityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, r.template(data), data.RandomInteger)
}
