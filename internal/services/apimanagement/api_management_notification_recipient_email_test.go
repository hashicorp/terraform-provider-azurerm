// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientemail"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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
	id, err := notificationrecipientemail.ParseRecipientEmailID(state.ID)
	if err != nil {
		return nil, err
	}

	notificationId := notificationrecipientemail.NewNotificationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.NotificationName)
	emails, err := client.ApiManagement.NotificationRecipientEmailClient.ListByNotificationComplete(ctx, notificationId)
	if err != nil {
		if !response.WasNotFound(emails.LatestHttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}
	}
	for _, existing := range emails.Items {
		if existing.Properties != nil && existing.Properties.Email != nil && *existing.Properties.Email == id.RecipientEmailName {
			return pointer.To(true), nil
		}
	}
	return pointer.To(false), nil
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
