// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDataShareDatasetKustoCluster_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_share_dataset_kusto_cluster", "test")
	r := DataShareDatasetKustoClusterResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue("azurerm_data_share_dataset_kusto_cluster.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_data_share_dataset_kusto_cluster.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_kusto_cluster.test", tfjsonpath.New("account_name"), tfjsonpath.New("share_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_kusto_cluster.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("share_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_data_share_dataset_kusto_cluster.test", tfjsonpath.New("share_name"), tfjsonpath.New("share_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
