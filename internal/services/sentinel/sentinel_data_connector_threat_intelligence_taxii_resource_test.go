// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type taxiiInfo struct {
	APIRootURL   string
	CollectionID string
	UserName     string
	Password     string
}

type DataConnectorThreatIntelligenceTAXIIResource struct {
	taxiiInfo    taxiiInfo
	taxiiInfoAlt taxiiInfo
}

func NewDataConnectorThreatIntelligenceTAXIIResource() DataConnectorThreatIntelligenceTAXIIResource {
	return DataConnectorThreatIntelligenceTAXIIResource{
		taxiiInfo: taxiiInfo{
			APIRootURL:   os.Getenv("ARM_TEST_TAXII_API_ROOT_URL"),
			CollectionID: os.Getenv("ARM_TEST_TAXII_COLLECTION_ID"),
			UserName:     os.Getenv("ARM_TEST_TAXII_USERNAME"),
			Password:     os.Getenv("ARM_TEST_TAXII_PASSWORD"),
		},
		taxiiInfoAlt: taxiiInfo{
			APIRootURL:   os.Getenv("ARM_TEST_TAXII_API_ROOT_URL_ALT"),
			CollectionID: os.Getenv("ARM_TEST_TAXII_COLLECTION_ID_ALT"),
			UserName:     os.Getenv("ARM_TEST_TAXII_USERNAME_ALT"),
			Password:     os.Getenv("ARM_TEST_TAXII_PASSWORD_ALT"),
		},
	}
}

func (r DataConnectorThreatIntelligenceTAXIIResource) preCheck(t *testing.T, forUpdate bool) {
	if r.taxiiInfo.APIRootURL == "" {
		t.Skipf(`"ARM_TEST_TAXII_API_ROOT_URL" not specified`)
	}
	if r.taxiiInfo.APIRootURL == "" {
		t.Skipf(`"ARM_TEST_TAXII_COLLECTION_ID" not specified`)
	}
	if forUpdate {
		if r.taxiiInfoAlt.APIRootURL == "" {
			t.Skipf(`"ARM_TEST_TAXII_API_ROOT_URL_ALT" not specified`)
		}
		if r.taxiiInfoAlt.CollectionID == "" {
			t.Skipf(`"ARM_TEST_TAXII_COLLECTION_ID_ALT" not specified`)
		}
	}
}

func TestAccDataConnectorThreatIntelligenceTAXII_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := NewDataConnectorThreatIntelligenceTAXIIResource()
	r.preCheck(t, false)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_name", "password"),
	})
}

func TestAccDataConnectorThreatIntelligenceTAXII_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := NewDataConnectorThreatIntelligenceTAXIIResource()
	r.preCheck(t, false)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_name", "password"),
	})
}

func TestAccDataConnectorThreatIntelligenceTAXII_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := NewDataConnectorThreatIntelligenceTAXIIResource()
	r.preCheck(t, true)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_name", "password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_name", "password"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("user_name", "password"),
	})
}

func TestAccDataConnectorThreatIntelligenceTAXII_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_data_connector_threat_intelligence_taxii", "test")
	r := NewDataConnectorThreatIntelligenceTAXIIResource()
	r.preCheck(t, false)

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

func (r DataConnectorThreatIntelligenceTAXIIResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Sentinel.DataConnectorsClient

	id, err := parse.DataConnectorID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r DataConnectorThreatIntelligenceTAXIIResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "test" {
  name                       = "acctestDC-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "test"
  api_root_url               = "%s"
  collection_id              = "%s"
  user_name                  = "%s"
  password                   = "%s"
}
`, template, data.RandomInteger, r.taxiiInfo.APIRootURL, r.taxiiInfo.CollectionID, r.taxiiInfo.UserName, r.taxiiInfo.Password)
}

func (r DataConnectorThreatIntelligenceTAXIIResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "test" {
  name                       = "acctestDC-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "test_update"
  api_root_url               = "%s"
  collection_id              = "%s"
  user_name                  = "%s"
  password                   = "%s"
  polling_frequency          = "OnceADay"
  lookback_date              = "1990-01-01T00:00:00Z"
}
`, template, data.RandomInteger, r.taxiiInfo.APIRootURL, r.taxiiInfo.CollectionID, r.taxiiInfo.UserName, r.taxiiInfo.Password)
}

func (r DataConnectorThreatIntelligenceTAXIIResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "test" {
  name                       = "acctestDC-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "test_update"
  api_root_url               = "%s"
  collection_id              = "%s"
  user_name                  = "%s"
  password                   = "%s"
  polling_frequency          = "OnceADay"
  lookback_date              = "1990-01-01T00:00:00Z"
}
`, template, data.RandomInteger, r.taxiiInfoAlt.APIRootURL, r.taxiiInfoAlt.CollectionID, r.taxiiInfoAlt.UserName, r.taxiiInfoAlt.Password)
}

func (r DataConnectorThreatIntelligenceTAXIIResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_data_connector_threat_intelligence_taxii" "import" {
  name                       = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.name
  log_analytics_workspace_id = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.log_analytics_workspace_id
  display_name               = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.display_name
  api_root_url               = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.api_root_url
  collection_id              = azurerm_sentinel_data_connector_threat_intelligence_taxii.test.collection_id
}
`, template)
}

func (r DataConnectorThreatIntelligenceTAXIIResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
