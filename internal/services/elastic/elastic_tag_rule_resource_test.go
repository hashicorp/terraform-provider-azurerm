package elastic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type TagRuleElasticMonitorResource struct{}

func TestAccElasticMonitorTagRule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_tag_rule", "test")
	r := TagRuleElasticMonitorResource{}
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

func TestAccElasticMonitorTagRule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_elastic_tag_rule", "test")
	r := TagRuleElasticMonitorResource{}
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

func (r TagRuleElasticMonitorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ElasticTagRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Elastic.TagRuleClient.Get(ctx, id.ResourceGroup, id.MonitorName, id.TagRuleName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Elastic Monitor rules %q (Resource Group %q): %+v", id.MonitorName, id.ResourceGroup, err)
	}

	if *resp.Properties.LogRules.SendActivityLogs == false && *resp.Properties.LogRules.SendSubscriptionLogs == false {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (r TagRuleElasticMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest-elastic-%d"
  location = "%s"
}
resource "azurerm_elastic_monitor" "test" {
  name                = "test-tf-elastic-tagrule-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  user_info {
    email_address = "ElasticTerraformTesting@mpliftrelastic20211117outlo.onmicrosoft.com"
  }
  sku {
    name = "staging_Monthly"
  }
}
`, data.RandomInteger%1000, data.Locations.Primary, data.RandomInteger%1000)
}

func (r TagRuleElasticMonitorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(
		`
	%s
	resource "azurerm_elastic_tag_rule" "test" {
		monitor_name = azurerm_elastic_monitor.test.name
		resource_group_name = azurerm_elastic_monitor.test.resource_group_name
		log_rules{
			send_subscription_logs = false
			send_activity_logs = true
			filtering_tag {
				name = "Test"
				value = "Terraform"
				action = "Include"
			}
		}
	}`, template)
}

func (r TagRuleElasticMonitorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(
		`
	%s
	resource "azurerm_elastic_tag_rule" "test" {
		monitor_name = azurerm_elastic_monitor.test.name
		resource_group_name = azurerm_elastic_monitor.test.resource_group_name
		log_rules{
			send_subscription_logs = true
			send_activity_logs = true
			filtering_tag {
				name = "Test"
				value = "Terraform"
				action = "Exclude"
			}
		}
	}`, template)
}
