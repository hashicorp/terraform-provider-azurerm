// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type AppServiceGithubTokenDataSource struct{}

func TestAccSourceControlGitHubTokenDataSource_basic(t *testing.T) {
	token := ""
	if token = os.Getenv("ARM_GITHUB_ACCESS_TOKEN"); token == "" {
		t.Skip("Skipping as `ARM_GITHUB_ACCESS_TOKEN` is not specified")
	}

	data := acceptance.BuildTestData(t, "data.azurerm_source_control_token", "test")
	r := AppServiceGithubTokenDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(token),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("token").Exists(),
			),
		},
	})
}

func (AppServiceGithubTokenDataSource) basic(token string) string {
	return fmt.Sprintf(`

%s

data azurerm_source_control_token test {
  type = azurerm_source_control_token.test.type
}
`, AppServiceGitHubTokenResource{}.basic(token))
}
