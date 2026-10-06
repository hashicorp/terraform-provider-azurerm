// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccDnsTxtRecord_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_txt_record", "test")
	r := DnsTxtRecordResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"record_type":         {},
		"dns_zone_name":       {},
		"name":                {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_dns_txt_record.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_dns_txt_record.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValue("azurerm_dns_txt_record.test", tfjsonpath.New("record_type"), knownvalue.StringExact("TXT")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dns_txt_record.test", tfjsonpath.New("dns_zone_name"), tfjsonpath.New("zone_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dns_txt_record.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_dns_txt_record.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
