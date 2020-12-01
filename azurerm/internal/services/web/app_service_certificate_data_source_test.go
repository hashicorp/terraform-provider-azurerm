package web_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AppServiceCertificateDataSource struct{}

func TestAccDataSourceAzureRMAppServiceCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_app_service_certificate", "test")

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: AppServiceCertificateDataSource{}.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("subject_name").Exists(),
				check.That(data.ResourceName).Key("issue_date").Exists(),
				check.That(data.ResourceName).Key("expiration_date").Exists(),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
	})
}

func (d AppServiceCertificateDataSource) basic(data acceptance.TestData) string {
	template := AppServiceCertificateResource{}.pfxNoPassword(data)
	return fmt.Sprintf(`
%s

data "azurerm_app_service_certificate" "test" {
  name                = azurerm_app_service_certificate.test.name
  resource_group_name = azurerm_app_service_certificate.test.resource_group_name
}
`, template)
}
