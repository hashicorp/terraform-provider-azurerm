// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudAppDynamicsApplicationPerformanceMonitoringResource struct{}

func TestAccSpringCloudAppDynamicsApplicationPerformanceMonitoring_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_dynamics_application_performance_monitoring", "test")
	r := SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("agent_account_name", "agent_account_access_key"),
	})
}

func TestAccSpringCloudAppDynamicsApplicationPerformanceMonitoring_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_dynamics_application_performance_monitoring", "test")
	r := SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
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

func TestAccSpringCloudAppDynamicsApplicationPerformanceMonitoring_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_dynamics_application_performance_monitoring", "test")
	r := SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("agent_account_name", "agent_account_access_key"),
	})
}

func TestAccSpringCloudAppDynamicsApplicationPerformanceMonitoring_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app_dynamics_application_performance_monitoring", "test")
	r := SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("agent_account_name", "agent_account_access_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("agent_account_name", "agent_account_access_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("agent_account_name", "agent_account_access_key"),
	})
}

func (r SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseApmID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.ApmsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_dynamics_application_performance_monitoring" "test" {
  name                     = "acctest-apm-%[2]d"
  spring_cloud_service_id  = azurerm_spring_cloud_service.test.id
  agent_account_name       = "test-agent-account-name"
  agent_account_access_key = "test-agent-account-access-key"
  controller_host_name     = "test-controller-host-name"
}
`, template, data.RandomInteger)
}

func (r SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_dynamics_application_performance_monitoring" "import" {
  name                     = azurerm_spring_cloud_app_dynamics_application_performance_monitoring.test.name
  spring_cloud_service_id  = azurerm_spring_cloud_app_dynamics_application_performance_monitoring.test.spring_cloud_service_id
  agent_account_name       = "test-agent-account-name"
  agent_account_access_key = "test-agent-account-access-key"
  controller_host_name     = "test-controller-host-name"
}
`, config)
}

func (r SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app_dynamics_application_performance_monitoring" "test" {
  name                     = "acctest-apm-%[2]d"
  spring_cloud_service_id  = azurerm_spring_cloud_service.test.id
  agent_account_name       = "updated-agent-account-name"
  agent_account_access_key = "updated-agent-account-access-key"
  controller_host_name     = "updated-controller-host-name"
  agent_application_name   = "test-agent-application-name"
  agent_tier_name          = "test-agent-tier-name"
  agent_node_name          = "test-agent-node-name"
  agent_unique_host_id     = "test-agent-unique-host-id"
  controller_ssl_enabled   = true
  controller_port          = 8080
  globally_enabled         = true
}
`, template, data.RandomInteger)
}
