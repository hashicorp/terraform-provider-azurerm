// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_notification_recipient_email

type ApiManagementNotificationRecipientEmailModel struct {
	ApiManagementId  string `tfschema:"api_management_id"`
	NotificationName string `tfschema:"notification_type"`
	Email            string `tfschema:"email"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementNotificationRecipientEmailModel{}
}
