// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridPartnerRegistrationDataSource struct{}

func TestAccDataSourceEventGridPartnerRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_partner_registration", "test")
	r := EventGridPartnerRegistrationDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("partner_registration_id").Exists(),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("test"),
			),
		},
	})
}

func (r EventGridPartnerRegistrationDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_partner_registration" "test" {
  name                = azurerm_eventgrid_partner_registration.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridPartnerRegistrationTestResource{}.basic(data))
}
