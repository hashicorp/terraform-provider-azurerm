// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccPrivateDnsCnameRecord_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_cname_record", "test")
	r := PrivateDnsCnameRecordResource{}

	checkedFields := map[string]struct{}{
		"name":                  {},
		"private_dns_zone_name": {},
		"record_type":           {},
		"resource_group_name":   {},
		"subscription_id":       {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_private_dns_cname_record.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_private_dns_cname_record.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_cname_record.test", tfjsonpath.New("private_dns_zone_name"), tfjsonpath.New("private_dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_cname_record.test", tfjsonpath.New("record_type"), tfjsonpath.New("id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_cname_record.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("private_dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_private_dns_cname_record.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("private_dns_zone_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
