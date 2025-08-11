// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type APIKeyDataSource struct{}

func TestAccAPIKeyDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_api_key", "test")
	secretText := uuid.NewString()
	endDateTime := getEndDateTime(3)
	r := APIKeyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data, secretText, endDateTime),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("end_date_time").HasValue(endDateTime),
				check.That(data.ResourceName).Key("hint").HasValue(secretText[:3]),
			),
		},
	})
}

func (APIKeyDataSource) complete(data acceptance.TestData, secretText, endDateTime string) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_api_key" "test" {
  name                = azurerm_nginx_api_key.test.name
  nginx_deployment_id = azurerm_nginx_api_key.test.nginx_deployment_id
}
`, APIKeyResource{}.complete(data, secretText, endDateTime))
}
