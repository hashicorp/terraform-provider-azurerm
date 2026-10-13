// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_notification_recipient_email

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientemail"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.Resource = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_api_management_notification_recipient_email"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return notificationrecipientemail.ValidateRecipientEmailID
}
