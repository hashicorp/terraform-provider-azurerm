package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementEmailTemplateResource struct {
}

func TestAccApiManagementEmailTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data, "basic"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementEmailTemplateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	templateName := apimanagement.TemplateName(id.Path["templates"])

	resp, err := clients.ApiManagement.EmailTemplateClient.Get(ctx, resourceGroup, serviceName, templateName)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Email Template (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementEmailTemplateResource) basic(data acceptance.TestData, testName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_email_template" "test" {
  template_name       = "ConfirmSignUpIdentityDefault"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  subject = "Please confirm your new customized $OrganizationName API account with this customized email"
  body = <<EOF
<!DOCTYPE html >
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Customized Letter Title</title>
  </head>
  <body>
    <table width="100%">
      <tr>
        <td>
          <p style="font-size:12pt;font-family:'Segoe UI'">Dear $DevFirstName $DevLastName,</p>
          <p style="font-size:12pt;font-family:'Segoe UI'"></p>
          <p style="font-size:12pt;font-family:'Segoe UI'">Thank you for joining the $OrganizationName API program! We host a growing number of cool APIs and strive to provide an awesome experience for API developers.</p>
          <p style="font-size:12pt;font-family:'Segoe UI'">This email is automatically created using a customized template witch is stored configuration as code.</p>
        </td>
      </tr>
    </table>
  </body>
</html>
EOF
}
`, r.template(data, testName), data.RandomInteger)
}

func (ApiManagementEmailTemplateResource) template(data acceptance.TestData, testName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d-%s"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}
`, data.RandomInteger, testName, data.Locations.Primary, data.RandomInteger, testName)
}
