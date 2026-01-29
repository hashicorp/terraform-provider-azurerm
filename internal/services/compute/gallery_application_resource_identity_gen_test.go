// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	customstatecheck "github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/statecheck"
)

func TestAccGalleryApplication_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application", "test")
	r := GalleryApplicationResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_gallery_application.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application.test", tfjsonpath.New("gallery_name"), tfjsonpath.New("gallery_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("gallery_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("gallery_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
