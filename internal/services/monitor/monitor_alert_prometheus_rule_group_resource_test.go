// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2023-03-01/prometheusrulegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AlertPrometheusRuleGroupTestResource struct{}

func TestAccAlertsManagementPrometheusRuleGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_prometheus_rule_group", "test")
	r := AlertPrometheusRuleGroupTestResource{}
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

func TestAccAlertsManagementPrometheusRuleGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_prometheus_rule_group", "test")
	r := AlertPrometheusRuleGroupTestResource{}
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

func TestAccAlertsManagementPrometheusRuleGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_prometheus_rule_group", "test")
	r := AlertPrometheusRuleGroupTestResource{}
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

func TestAccAlertsManagementPrometheusRuleGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_alert_prometheus_rule_group", "test")
	r := AlertPrometheusRuleGroupTestResource{}
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

func (r AlertPrometheusRuleGroupTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := prometheusrulegroups.ParsePrometheusRuleGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Monitor.AlertPrometheusRuleGroupClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r AlertPrometheusRuleGroupTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_monitor_workspace" "test" {
  name                = "acctest-amw-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[2]s"
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r AlertPrometheusRuleGroupTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_monitor_alert_prometheus_rule_group" "test" {
  name                = "acctest-amprg-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  scopes              = [azurerm_monitor_workspace.test.id]
  rule {
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    record     = "job_type:billing_jobs_duration_seconds:99p5m"
    labels = {
      team = "prod"
    }
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AlertPrometheusRuleGroupTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_alert_prometheus_rule_group" "import" {
  name                = azurerm_monitor_alert_prometheus_rule_group.test.name
  resource_group_name = azurerm_monitor_alert_prometheus_rule_group.test.resource_group_name
  location            = azurerm_monitor_alert_prometheus_rule_group.test.location
  scopes              = azurerm_monitor_alert_prometheus_rule_group.test.scopes

  rule {
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    record     = "job_type:billing_jobs_duration_seconds:99p5m"
    labels = {
      team = "prod"
    }
  }
}
`, config)
}

func (r AlertPrometheusRuleGroupTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_monitor_alert_prometheus_rule_group" "test" {
  name                = "acctest-amprg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  cluster_name        = azurerm_kubernetes_cluster.test.name
  description         = "This is the description of the following rule group"
  rule_group_enabled  = false
  interval            = "PT10M"
  scopes              = [azurerm_monitor_workspace.test.id]
  rule {
    enabled    = false
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    record     = "job_type:billing_jobs_duration_seconds:99p5m"
    labels = {
      team = "prod"
    }
  }
  rule {
    alert      = "Billing_Processing_Very_Slow"
    enabled    = true
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    for        = "PT5M"
    severity   = 0
    action {
      action_group_id = azurerm_monitor_action_group.test.id
    }
    alert_resolution {
      auto_resolved   = true
      time_to_resolve = "PT10M"
    }
    annotations = {
      annotationName = "annotationValue"
    }
    labels = {
      team = "prod"
    }
  }
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AlertPrometheusRuleGroupTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  default_node_pool {
    name                    = "default"
    node_count              = 1
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_cluster" "test2" {
  name                = "acctestaks2%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks2%[2]d"

  default_node_pool {
    name                    = "default"
    node_count              = 2
    vm_size                 = "Standard_DS2_v2"
    host_encryption_enabled = true
    upgrade_settings {
      max_surge = "10%%"
    }
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag2"
}

resource "azurerm_monitor_alert_prometheus_rule_group" "test" {
  name                = "acctest-amprg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"
  cluster_name        = azurerm_kubernetes_cluster.test2.name
  description         = "This is the description of the following rule group2"
  rule_group_enabled  = true
  interval            = "PT10M"
  scopes              = [azurerm_monitor_workspace.test.id]

  rule {
    enabled    = true
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    record     = "job_type:billing_jobs_duration_seconds:99p6m"
    labels = {
      team2 = "prod2"
    }
  }

  rule {
    alert      = "Billing_Processing_Very_Slow2"
    enabled    = false
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    for        = "PT4M"
    severity   = 1
    action {
      action_group_id = azurerm_monitor_action_group.test2.id
    }
    action {
      action_group_id = azurerm_monitor_action_group.test.id
    }
    alert_resolution {
      auto_resolved   = false
      time_to_resolve = "PT9M"
    }
    annotations = {
      annotationName2 = "annotationValue2"
    }
    labels = {
      team2 = "prod2"
    }
  }
  tags = {
    key2 = "value2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
