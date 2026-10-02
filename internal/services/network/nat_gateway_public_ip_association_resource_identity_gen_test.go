// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccNatGatewayPublicIpAssociation_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_nat_gateway_public_ip_association", "test")
	r := NatGatewayPublicIpAssociationResource{}

	checkedFields := map[string]struct{}{
		"nat_gateway_id":       {},
		"public_ip_address_id": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_nat_gateway_public_ip_association.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_nat_gateway_public_ip_association.test", tfjsonpath.New("nat_gateway_id"), tfjsonpath.New("nat_gateway_id")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_nat_gateway_public_ip_association.test", tfjsonpath.New("public_ip_address_id"), tfjsonpath.New("public_ip_address_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
