// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2023-06-01/monitorsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticsearchResource struct{}

func TestAccElasticsearch_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_cloud_deployment_id").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_sso_default_url").Exists(),
				check.That(data.ResourceName).Key("elastic_cloud_user_id").Exists(),
				check.That(data.ResourceName).Key("elasticsearch_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_sso_uri").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticsearch_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
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

func TestAccElasticsearch_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
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

func TestAccElasticsearch_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticsearch_logs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// this proves that we don't need to destroy the `logs` block separately
			Config: r.logs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticsearch_logsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_elasticsearch", "test")
	r := ElasticsearchResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// create with it
			Config: r.logs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// update it
			Config: r.logsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// remove just the `logs` block
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ElasticsearchResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitorsresource.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Elastic.MonitorClient.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ElasticsearchResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"

  lifecycle {
    ignore_changes = [logs]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticsearchResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_cloud_elasticsearch" "import" {
  name                        = azurerm_elastic_cloud_elasticsearch.test.name
  resource_group_name         = azurerm_elastic_cloud_elasticsearch.test.resource_group_name
  location                    = azurerm_elastic_cloud_elasticsearch.test.location
  sku_name                    = azurerm_elastic_cloud_elasticsearch.test.sku_name
  elastic_cloud_email_address = azurerm_elastic_cloud_elasticsearch.test.elastic_cloud_email_address
}
`, r.basic(data))
}

func (r ElasticsearchResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticsearchResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
  monitoring_enabled          = false

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticsearchResource) logs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"

  logs {
    filtering_tag {
      action = "Include"
      name   = "TerraformAccTest"
      value  = "RandomValue%[1]d"
    }

    # NOTE: these are intentionally not set to true here for testing purposes
    send_activity_logs     = false
    send_azuread_logs      = false
    send_subscription_logs = false
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ElasticsearchResource) logsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-estc%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"

  logs {
    filtering_tag {
      action = "Include"
      name   = "TerraformAccTest"
      value  = "UpdatedValue-%[1]d"
    }

    # NOTE: these are intentionally not set to true here for testing purposes
    send_activity_logs     = false
    send_azuread_logs      = false
    send_subscription_logs = false
  }

  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
