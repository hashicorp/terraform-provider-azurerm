// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IoTCentralApplicationNetworkRuleSetResource struct{}

func TestAccIoTCentralApplicationNetworkRuleSet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("apply_to_device").HasValue("true"),
				check.That(data.ResourceName).Key("default_action").HasValue("Deny"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplicationNetworkRuleSet_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

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

func TestAccIoTCentralApplicationNetworkRuleSet_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

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
	})
}

func TestAccIoTCentralApplicationNetworkRuleSet_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

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

func TestAccIoTCentralApplicationNetworkRuleSet_applyToDevice(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.applyToDevice(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.applyToDeviceUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplicationNetworkRuleSet_defaultAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultAction(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.defaultActionUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplicationNetworkRuleSet_IPRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.IPRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.IPRuleUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplicationNetworkRuleSet_updateIoTCentralApplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application_network_rule_set", "test")
	r := IoTCentralApplicationNetworkRuleSetResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ioTCentralApplicationPublicNetworkAccessEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ioTCentralApplicationPublicNetworkAccessDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (IoTCentralApplicationNetworkRuleSetResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := apps.ParseIotAppID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTCentral.AppsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.NetworkRuleSets != nil), nil
}

func (r IoTCentralApplicationNetworkRuleSetResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  apply_to_device = false
  default_action  = "Allow"

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.1.0/24"
  }

  ip_rule {
    name    = "rule2"
    ip_mask = "10.1.1.0/24"
  }
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "import" {
  iotcentral_application_id = azurerm_iotcentral_application_network_rule_set.test.iotcentral_application_id
}
`, template)
}

func (r IoTCentralApplicationNetworkRuleSetResource) applyToDevice(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  apply_to_device = false
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) applyToDeviceUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  apply_to_device = true
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) defaultAction(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  default_action = "Allow"
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) defaultActionUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  default_action = "Deny"
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) IPRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.1.0/24"
  }
}
`, r.template(data))
}

func (r IoTCentralApplicationNetworkRuleSetResource) IPRuleUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.3.0/24"
  }

  ip_rule {
    name    = "rule2"
    ip_mask = "10.1.1.0/24"
  }
}
`, r.template(data))
}

func (IoTCentralApplicationNetworkRuleSetResource) ioTCentralApplicationPublicNetworkAccessEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"

  public_network_access_enabled = true
}

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.1.0/24"
  }

  ip_rule {
    name    = "rule2"
    ip_mask = "10.1.1.0/24"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationNetworkRuleSetResource) ioTCentralApplicationPublicNetworkAccessDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"

  public_network_access_enabled = false
}

resource "azurerm_iotcentral_application_network_rule_set" "test" {
  iotcentral_application_id = azurerm_iotcentral_application.test.id

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.1.0/24"
  }

  ip_rule {
    name    = "rule2"
    ip_mask = "10.1.1.0/24"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationNetworkRuleSetResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "subdomain-%[1]d"
}
`, data.RandomInteger, data.Locations.Primary)
}
