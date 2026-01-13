// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datadog_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccDatadogMonitorTagRule_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_datadog_monitor_tag_rule", "test")
	r := DatadogMonitorTagRuleResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue("azurerm_datadog_monitor_tag_rule.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_datadog_monitor_tag_rule.test", tfjsonpath.New("name"), tfjsonpath.New("tag_rule_name")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_datadog_monitor_tag_rule.test", tfjsonpath.New("monitor_name"), tfjsonpath.New("datadog_monitor_id")),
					customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_datadog_monitor_tag_rule.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("datadog_monitor_id")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(),
			data.ImportBlockWithIDStep(),
		},
	})
}
