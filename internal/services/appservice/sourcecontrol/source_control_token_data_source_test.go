package sourcecontrol_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
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
