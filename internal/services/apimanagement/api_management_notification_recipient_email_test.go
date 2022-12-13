package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiManagementNotificationRecipientEmailResource struct{}

func TestAccApiManagementNotificationRecipientEmail_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_notification_recipient_email", "test")
	r := ApiManagementNotificationRecipientEmailResource{}

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

func TestAccApiManagementNotificationRecipientEmail_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_notification_recipient_email", "test")
	r := ApiManagementNotificationRecipientEmailResource{}

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

func (ApiManagementNotificationRecipientEmailResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.NotificationRecipientEmailID(state.ID)
	if err != nil {
		return nil, err
	}

	emails, err := client.ApiManagement.NotificationRecipientEmailClient.ListByNotification(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName))
	if err != nil {
		if !utils.ResponseWasNotFound(emails.Response) {
			return nil, fmt.Errorf("retrieving Api Management Notification Recipient Email %q (Resource Group %q): %+v", id.RecipientEmailName, id.ResourceGroup, err)
		}
	}
	if emails.Value != nil {
		for _, existing := range *emails.Value {
			if existing.RecipientEmailContractProperties != nil && existing.RecipientEmailContractProperties.Email != nil && *existing.RecipientEmailContractProperties.Email == id.RecipientEmailName {
				return utils.Bool(true), nil
			}
		}
	}
	return utils.Bool(false), nil
}

func (r ApiManagementNotificationRecipientEmailResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_api_management_notification_recipient_email" "test" {
  api_management_id = azurerm_api_management.test.id
  notification_type = "AccountClosedPublisher"
  email             = "foo@bar.com"
}
`, ApiManagementResource{}.basic(data))
}

func (r ApiManagementNotificationRecipientEmailResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_notification_recipient_email" "import" {
  api_management_id = azurerm_api_management.test.id
  notification_type = azurerm_api_management_notification_recipient_email.test.notification_type
  email             = azurerm_api_management_notification_recipient_email.test.email
}
`, r.basic(data))
}
