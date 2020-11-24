package authorization_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ClientConfigDataSource struct{}

func TestAccDataSourceAzureRMClientConfig_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_client_config", "current")
	clientId := os.Getenv("ARM_CLIENT_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	objectIdRegex := regexp.MustCompile("^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: ClientConfigDataSource{}.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("client_id").HasValue(clientId),
				check.That(data.ResourceName).Key("tenant_id").HasValue(tenantId),
				check.That(data.ResourceName).Key("subscription_id").HasValue(subscriptionId),
				check.That(data.ResourceName).Key("object_id").MatchesRegex(objectIdRegex),
			),
		},
	})
}

func (d ClientConfigDataSource) basic() string {
	return `
data "azurerm_client_config" "current" { }
`
}
