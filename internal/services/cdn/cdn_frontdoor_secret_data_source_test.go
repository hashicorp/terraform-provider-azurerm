// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type CdnFrontdoorSecretResourceDataSource struct {
	DoNotRunFrontDoorCustomDomainTests string
}

// NOTE: This is currently not testable due to the cert requirements of the service
func TestAccCdnFrontDoorSecretDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_frontdoor_secret", "test")
	r := CdnFrontdoorSecretResource{os.Getenv("ARM_TEST_DO_NOT_RUN_CDN_FRONT_DOOR_CUSTOM_DOMAIN")}
	r.preCheck(t)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("cdn_frontdoor_profile_id").MatchesOtherKey(check.That("azurerm_cdn_frontdoor_profile.test").Key("id")),
			),
		},
		data.ImportStep(),
	})
}
