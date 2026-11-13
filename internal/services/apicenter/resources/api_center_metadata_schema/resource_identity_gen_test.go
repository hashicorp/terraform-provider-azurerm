// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccApiCenterMetadataSchema_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_metadata_schema", "test")
	r := ApiCenterMetadataSchemaResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"service_name":        {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_api_center_metadata_schema.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_api_center_metadata_schema.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_metadata_schema.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_metadata_schema.test", tfjsonpath.New("service_name"), tfjsonpath.New("service_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_api_center_metadata_schema.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("service_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
