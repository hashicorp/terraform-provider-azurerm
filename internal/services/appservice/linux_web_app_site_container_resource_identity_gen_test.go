// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
	"testing"
)

func TestAccLinuxWebAppSiteContainer_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_linux_web_app_site_container", "test")
	r := LinuxWebAppSiteContainerResource{}

	checkedFields := map[string]struct{}{
		"name":                {},
		"resource_group_name": {},
		"site_name":           {},
		"subscription_id":     {},
	}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				customstatecheck.ExpectAllIdentityFieldsAreChecked("azurerm_linux_web_app_site_container.test", checkedFields),
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_linux_web_app_site_container.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_linux_web_app_site_container.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("linux_web_app_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_linux_web_app_site_container.test", tfjsonpath.New("site_name"), tfjsonpath.New("linux_web_app_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_linux_web_app_site_container.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("linux_web_app_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
