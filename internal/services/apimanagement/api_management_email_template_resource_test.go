// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/emailtemplates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementEmailTemplateResource struct{}

func TestAccApiManagementEmailTemplate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementEmailTemplate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccApiManagementEmailTemplate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_email_template", "test")
	r := ApiManagementEmailTemplateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subject").HasValue("Please confirm your new customized $OrganizationName API account with this customized email"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("subject").HasValue("Please confirm your new customized $OrganizationName API account with this customized and updated email"),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementEmailTemplateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := emailtemplates.ParseTemplateIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	templateName := emailtemplates.TemplateName(azure.TitleCase(string(id.TemplateName)))
	newId := emailtemplates.NewTemplateID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, templateName)

	_, err = clients.ApiManagement.EmailTemplatesClient.EmailTemplateGet(ctx, newId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", newId, err)
	}

	return pointer.To(true), nil
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
