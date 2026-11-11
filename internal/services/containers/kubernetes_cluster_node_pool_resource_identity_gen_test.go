// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccKubernetesClusterNodePool_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_cluster_node_pool", "test")
	r := KubernetesClusterNodePoolResource{}

	checkedFields := map[string]struct{}{
		"name":                 {},
		"managed_cluster_name": {},
		"resource_group_name":  {},
		"subscription_id":      {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.manualScaleConfig(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_kubernetes_cluster_node_pool.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_kubernetes_cluster_node_pool.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_cluster_node_pool.test", tfjsonpath.New("managed_cluster_name"), tfjsonpath.New("kubernetes_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_cluster_node_pool.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("kubernetes_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_kubernetes_cluster_node_pool.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("kubernetes_cluster_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
