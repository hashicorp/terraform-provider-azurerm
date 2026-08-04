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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
				check.That(data.ResourceName).Key("sku_name").HasValue("ess-consumption-2024_Monthly"),
				check.That(data.ResourceName).Key("project_type").HasValue("Elasticsearch"),
				check.That(data.ResourceName).Key("configuration_type").HasValue("GeneralPurpose"),
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

func TestAccElasticCloudServerless_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_serverless", "test")
	r := ElasticCloudServerlessResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("LiveTest"),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, plancheck.ResourceActionUpdate),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.Environment").HasValue("LiveTest"),
				check.That(data.ResourceName).Key("tags.TagRevision").HasValue("after"),
				check.That(data.ResourceName).Key("tags.UpdateMode").HasValue("InPlace"),
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
	return r.config(data, "ess-consumption-2024_Monthly", "elastic-serverless-search", "Elasticsearch", "GeneralPurpose", "ec-azure-pp", "n7ja87drquhy", "")
}

func (r ElasticCloudServerlessResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_elastic_cloud_serverless" "import" {
  name                        = azurerm_elastic_cloud_serverless.test.name
  resource_group_name         = azurerm_elastic_cloud_serverless.test.resource_group_name
  location                    = azurerm_elastic_cloud_serverless.test.location
  kind                        = azurerm_elastic_cloud_serverless.test.kind
  sku_name                    = azurerm_elastic_cloud_serverless.test.sku_name
  project_type                = azurerm_elastic_cloud_serverless.test.project_type
  configuration_type          = azurerm_elastic_cloud_serverless.test.configuration_type
  offer_id                    = azurerm_elastic_cloud_serverless.test.offer_id
  term_id                     = azurerm_elastic_cloud_serverless.test.term_id
  elastic_cloud_email_address = azurerm_elastic_cloud_serverless.test.elastic_cloud_email_address
}
`, r.basic(data))
}

func (r ElasticCloudServerlessResource) complete(data acceptance.TestData) string {
	return r.config(data, "ess-consumption-2024_Monthly", "elastic-serverless-search", "Elasticsearch", "GeneralPurpose", "ec-azure-pp", "n7ja87drquhy", `
	generate_api_key    = false
	monitoring_enabled = true
	plan_id             = "ess-consumption-2024"
	publisher_id        = "elastic"

  tags = {
    Environment = "LiveTest"
    TagRevision = "before"
  }`)
}

func (r ElasticCloudServerlessResource) completeUpdated(data acceptance.TestData) string {
	return r.config(data, "ess-consumption-2024_Monthly", "elastic-serverless-search", "Elasticsearch", "GeneralPurpose", "ec-azure-pp", "n7ja87drquhy", `
	generate_api_key    = false
	monitoring_enabled = true
	plan_id             = "ess-consumption-2024"
	publisher_id        = "elastic"

  tags = {
    Environment = "LiveTest"
    TagRevision = "after"
    UpdateMode   = "InPlace"
  }`)
}

func (r ElasticCloudServerlessResource) config(data acceptance.TestData, skuName, kind, projectType, configurationType, offerId, termId, tags string) string {
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
  kind                        = %[2]q
  sku_name                    = %[3]q
  project_type                = %[4]q
  configuration_type          = %[5]q
  offer_id                    = %[6]q
  term_id                     = %[7]q
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
%[8]s
}
`, data.RandomInteger, kind, skuName, projectType, configurationType, offerId, termId, tags)
}
