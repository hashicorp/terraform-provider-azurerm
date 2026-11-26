// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redhatopenshift_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

// A bug in the HCP service currently prevents proper destruction of node pool resources
// because it won't allow deletion of a cluster's last node pool once at least one has been created.
// Once the bug is fixed, we can replace this test file with a generated one.
func TestAccRedhatHcpOpenshiftNodePool_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_node_pool", "test")
	r := RedhatHcpOpenshiftNodePoolResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":             {},
		"hcp_open_shift_cluster_name": {},
		"name":                        {},
		"resource_group_name":         {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_redhat_hcp_openshift_node_pool.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_redhat_hcp_openshift_node_pool.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_node_pool.test", tfjsonpath.New("hcp_open_shift_cluster_name"), tfjsonpath.New("hcp_openshift_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_node_pool.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("hcp_openshift_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_node_pool.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("hcp_openshift_cluster_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(true),
		data.ImportBlockWithIDStep(true),
		{
			Config: r.removeFromState(RedhatHcpOpenshiftClusterResource{}.basic(data)),
		},
	}, false)
}
