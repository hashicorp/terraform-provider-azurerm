// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccMssqlElasticPool_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elastic_pool", "test")
	r := MssqlElasticPoolResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"resource_group_name": {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basicDTU(data, ""),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_mssql_elastic_pool.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_mssql_elastic_pool.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_elastic_pool.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_elastic_pool.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
