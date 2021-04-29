package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ApiManagementEmailTemplateResource struct {
}

func TestAccApiManagementEmailTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementEmailTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementEmailTemplate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subject").HasValue("Please confirm your new customized $OrganizationName API account with this customized email"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subject").HasValue("Please confirm your new customized $OrganizationName API account with this customized and updated email"),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementEmailTemplateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.EmailTemplateID(state.ID)
	if err != nil {
		return nil, err
	}

	_, err = clients.ApiManagement.EmailTemplateClient.Get(ctx, id.ResourceGroup, id.ServiceName, apimanagement.TemplateName(id.TemplateName))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ApiManagementEmailTemplateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_email_template" "test" {
  template_name       = "ConfirmSignUpIdentityDefault"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  subject             = "Please confirm your new customized $OrganizationName API account with this customized email"
  body                = <<EOF
<!DOCTYPE html >
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Customized Letter Title</title>
  </head>
  <body>
    <table width="100%%">
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
`, r.template(data))
}

func (r ApiManagementEmailTemplateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_email_template" "import" {
  template_name       = azurerm_api_management_email_template.test.template_name
  api_management_name = azurerm_api_management_email_template.test.api_management_name
  resource_group_name = azurerm_api_management_email_template.test.resource_group_name
  subject             = azurerm_api_management_email_template.test.subject
  body                = azurerm_api_management_email_template.test.body
}
`, r.basic(data))
}

func (r ApiManagementEmailTemplateResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_email_template" "test" {
  template_name       = "ConfirmSignUpIdentityDefault"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  subject             = "Please confirm your new customized $OrganizationName API account with this customized and updated email"
  body                = <<EOF
<!DOCTYPE html >
<html>
  <head>
    <meta charset="UTF-8" />
    <title>Customized Letter Title</title>
  </head>
  <body>
    <table width="100%%">
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
`, r.template(data))
}

func (ApiManagementEmailTemplateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
