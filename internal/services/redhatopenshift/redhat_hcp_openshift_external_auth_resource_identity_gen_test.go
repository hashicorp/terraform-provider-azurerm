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

func TestAccRedhatHcpOpenshiftExternalAuth_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redhat_hcp_openshift_external_auth", "test")
	r := RedhatHcpOpenshiftExternalAuthResource{}

	checkedFields := map[string]struct{}{
		"name":                        {},
		"hcp_open_shift_cluster_name": {},
		"resource_group_name":         {},
		"subscription_id":             {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_redhat_hcp_openshift_external_auth.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_redhat_hcp_openshift_external_auth.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_external_auth.test", tfjsonpath.New("hcp_open_shift_cluster_name"), tfjsonpath.New("hcp_openshift_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_external_auth.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("hcp_openshift_cluster_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_redhat_hcp_openshift_external_auth.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("hcp_openshift_cluster_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
