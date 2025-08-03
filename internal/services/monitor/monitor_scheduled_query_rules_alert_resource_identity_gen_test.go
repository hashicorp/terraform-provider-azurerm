// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccMonitorScheduledQueryRulesAlert_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesAlertResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0-rc2"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.AlertingActionConfigComplete(data),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentityValue("azurerm_monitor_scheduled_query_rules_alert.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_monitor_scheduled_query_rules_alert.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
					statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_monitor_scheduled_query_rules_alert.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				},
			},
			data.ImportBlockWithResourceIdentityStep(),
			data.ImportBlockWithIDStep(),
		},
	})
}
