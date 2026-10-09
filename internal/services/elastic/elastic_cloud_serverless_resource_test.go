// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elastic/2025-06-01/elasticmonitorresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ElasticCloudServerlessResource struct{}

func TestAccElasticCloudServerless_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_serverless", "test")
	r := ElasticCloudServerlessResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_cloud_deployment_id").Exists(),
				check.That(data.ResourceName).Key("elasticsearch_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_service_url").Exists(),
				check.That(data.ResourceName).Key("kibana_sso_uri").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccElasticCloudServerless_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_serverless", "test")
	r := ElasticCloudServerlessResource{}
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

func TestAccElasticCloudServerless_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_serverless", "test")
	r := ElasticCloudServerlessResource{}
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

func (r ElasticCloudServerlessResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := elasticmonitorresources.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Elastic.ServerlessMonitorClient.MonitorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r ElasticCloudServerlessResource) basic(data acceptance.TestData) string {
	return r.config(data, "")
}

func (r ElasticCloudServerlessResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_cloud_serverless" "import" {
  name                        = azurerm_elastic_cloud_serverless.test.name
  resource_group_name         = azurerm_elastic_cloud_serverless.test.resource_group_name
  location                    = azurerm_elastic_cloud_serverless.test.location
  configuration_type          = azurerm_elastic_cloud_serverless.test.configuration_type
  elastic_cloud_email_address = azurerm_elastic_cloud_serverless.test.elastic_cloud_email_address
  kind                        = azurerm_elastic_cloud_serverless.test.kind
  offer_id                    = azurerm_elastic_cloud_serverless.test.offer_id
  project_type                = azurerm_elastic_cloud_serverless.test.project_type
  sku                         = azurerm_elastic_cloud_serverless.test.sku
  term_id                     = azurerm_elastic_cloud_serverless.test.term_id
}
`, r.basic(data))
}

func (r ElasticCloudServerlessResource) complete(data acceptance.TestData) string {
	return r.config(data, `
  generate_api_key    = false
  monitoring_enabled = true
  plan_id            = "ess-consumption-2024"
  publisher_id       = "elastic"

  tags = {
    Environment = "LiveTest"
    TagRevision = "before"
  }`)
}

func (r ElasticCloudServerlessResource) config(data acceptance.TestData, optionalConfig string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-serverless-%[1]d"
  location = "eastus2euap"
}

resource "azurerm_elastic_cloud_serverless" "test" {
  name                        = "acctest-es-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  configuration_type          = "GeneralPurpose"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
  kind                        = "elastic-serverless-search"
  offer_id                    = "ec-azure-pp"
  project_type                = "Elasticsearch"
  sku                         = "ess-consumption-2024_Monthly"
  term_id                     = "n7ja87drquhy"
%[2]s
}
`, data.RandomInteger, optionalConfig)
}
