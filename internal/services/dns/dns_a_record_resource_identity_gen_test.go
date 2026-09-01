// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDnsARecord_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_a_record", "test")
	r := DnsARecordResource{}

	checkedFields := map[string]struct{}{
		"dns_zone_name":       {},
		"name":                {},
		"record_type":         {},
		"resource_group_name": {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_dns_a_record.test", checkedFields),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dns_a_record.test", tfjsonpath.New("dns_zone_name"), tfjsonpath.New("dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dns_a_record.test", tfjsonpath.New("name"), tfjsonpath.New("dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dns_a_record.test", tfjsonpath.New("record_type"), tfjsonpath.New("dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dns_a_record.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("dns_zone_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_dns_a_record.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("dns_zone_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
