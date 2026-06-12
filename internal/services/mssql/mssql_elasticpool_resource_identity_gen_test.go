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

func TestAccMssqlElasticpool_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_elasticpool", "test")
	r := MssqlElasticpoolResource{}

	checkedFields := map[string]struct{}{
		"subscription_id":     {},
		"name":                {},
		"resource_group_name": {},
		"server_name":         {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_mssql_elasticpool.test", checkedFields),
				statecheck.ExpectIdentityValue("azurerm_mssql_elasticpool.test", tfjsonpath.New("subscription_id"), knownvalue.StringExact(data.Subscriptions.Primary)),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_elasticpool.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_elasticpool.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("resource_group_name")),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_mssql_elasticpool.test", tfjsonpath.New("server_name"), tfjsonpath.New("server_name")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
