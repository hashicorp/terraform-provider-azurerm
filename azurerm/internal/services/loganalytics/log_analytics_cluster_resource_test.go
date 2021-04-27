package loganalytics_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogAnalyticsClusterResource struct {
}

func TestAccLogAnalyticsCluster_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	r := LogAnalyticsClusterResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsCluster_resize(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	r := LogAnalyticsClusterResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.resize(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogAnalyticsCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_log_analytics_cluster", "test")
	r := LogAnalyticsClusterResource{}

	if os.Getenv("ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS") == "" {
		t.Skip("Skipping as ARM_RUN_TEST_LOG_ANALYTICS_CLUSTERS is not specified")
		return
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t LogAnalyticsClusterResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.LogAnalyticsClusterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.LogAnalytics.ClusterClient.Get(ctx, id.ResourceGroup, id.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("readingLog Analytics Cluster (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (LogAnalyticsClusterResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-la-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogAnalyticsClusterResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsClusterResource) resize(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "test" {
  name                = "acctest-LA-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size_gb             = 1100

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r LogAnalyticsClusterResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_log_analytics_cluster" "import" {
  name                = azurerm_log_analytics_cluster.test.name
  resource_group_name = azurerm_log_analytics_cluster.test.resource_group_name
  location            = azurerm_log_analytics_cluster.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data))
}
