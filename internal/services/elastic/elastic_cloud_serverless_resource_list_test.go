// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccElasticCloudServerless_list_basic(t *testing.T) {
	r := ElasticCloudServerlessResource{}
	listResourceAddress := "azurerm_elastic_cloud_serverless.list"
	data := acceptance.BuildTestData(t, "azurerm_elastic_cloud_serverless", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r ElasticCloudServerlessResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-elastic-serverless-list-%[1]d"
  location = "eastus2euap"
}

resource "azurerm_elastic_cloud_serverless" "test" {
  count = 2

  name                        = "acctest-es-srv-${count.index}-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  kind                        = "elastic-serverless-search"
  sku_name                    = "ess-consumption-2024_Monthly"
  project_type                = "Elasticsearch"
  configuration_type          = "GeneralPurpose"
  offer_id                    = "ec-azure-pp"
  term_id                     = "n7ja87drquhy"
  elastic_cloud_email_address = "terraform-acctest@hashicorp.com"
}
`, data.RandomInteger)
}

func (r ElasticCloudServerlessResource) basicQuery() string {
	return `
list "azurerm_elastic_cloud_serverless" "list" {
  provider = azurerm
  config {}
}
`
}

func (r ElasticCloudServerlessResource) basicQueryByResourceGroupName(data acceptance.TestData) string {
	return fmt.Sprintf(`
list "azurerm_elastic_cloud_serverless" "list" {
  provider = azurerm
  config {
    resource_group_name = "acctestrg-elastic-serverless-list-%[1]d"
  }
}
`, data.RandomInteger)
}
