// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementNotificationRecipientUserResource struct{}

func TestAccApiManagementNotificationRecipientUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_notification_recipient_user", "test")
	r := ApiManagementNotificationRecipientUserResource{}

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

func TestAccApiManagementNotificationRecipientUser_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_notification_recipient_user", "test")
	r := ApiManagementNotificationRecipientUserResource{}

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

func (ApiManagementNotificationRecipientUserResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := notificationrecipientuser.ParseRecipientUserID(state.ID)
	if err != nil {
		return nil, err
	}

	notificationId := notificationrecipientuser.NewNotificationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.NotificationName)
	users, err := client.ApiManagement.NotificationRecipientUserClient.ListByNotificationComplete(ctx, notificationId)
	if err != nil {
		if !response.WasNotFound(users.LatestHttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}
	}
	for _, existing := range users.Items {
		if existing.Name != nil && *existing.Name == id.UserId {
			return pointer.To(true), nil
		}
	}
	return pointer.To(false), nil
}

func (r ApiManagementNotificationRecipientUserResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_user" "test" {
  user_id             = "123"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Example"
  last_name           = "User"
  email               = "foo@bar.com"
  state               = "active"
}

resource "azurerm_api_management_notification_recipient_user" "test" {
  api_management_id = azurerm_api_management.test.id
  notification_type = "AccountClosedPublisher"
  user_id           = azurerm_api_management_user.test.user_id
}
`, ApiManagementResource{}.basic(data))
}

func (r ApiManagementNotificationRecipientUserResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_notification_recipient_user" "import" {
  api_management_id = azurerm_api_management.test.id
  notification_type = azurerm_api_management_notification_recipient_user.test.notification_type
  user_id           = azurerm_api_management_notification_recipient_user.test.user_id
}
`, r.basic(data))
}
