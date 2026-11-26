// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccSubnetNetworkSecurityGroupAssociation_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subnet_network_security_group_association", "test")
	r := SubnetNetworkSecurityGroupAssociationResource{}

	checkedFields := map[string]struct{}{
		"name":                 {},
		"resource_group_name":  {},
		"subscription_id":      {},
		"virtual_network_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_subnet_network_security_group_association.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_subnet_network_security_group_association.test", tfjsonpath.New("name"), tfjsonpath.New("subnet_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_subnet_network_security_group_association.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("subnet_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_subnet_network_security_group_association.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("subnet_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_subnet_network_security_group_association.test", tfjsonpath.New("virtual_network_name"), tfjsonpath.New("subnet_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
