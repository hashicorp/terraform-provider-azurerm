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

func TestAccGalleryApplicationVersion_resourceIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_gallery_application_version", "test")
	r := GalleryApplicationVersionResource{}

	data.ResourceIdentityTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValueMatchesStateAtPath("azurerm_gallery_application_version.test", tfjsonpath.New("name"), tfjsonpath.New("name")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application_version.test", tfjsonpath.New("application_name"), tfjsonpath.New("gallery_application_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application_version.test", tfjsonpath.New("gallery_name"), tfjsonpath.New("gallery_application_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application_version.test", tfjsonpath.New("resource_group_name"), tfjsonpath.New("gallery_application_id")),
				customstatecheck.ExpectStateContainsIdentityValueAtPath("azurerm_gallery_application_version.test", tfjsonpath.New("subscription_id"), tfjsonpath.New("gallery_application_id")),
			},
		},
		data.ImportBlockWithResourceIdentityStep(false),
		data.ImportBlockWithIDStep(false),
	}, false)
}
