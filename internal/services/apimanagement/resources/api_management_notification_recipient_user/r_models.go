// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_notification_recipient_user

type ApiManagementNotificationRecipientUserModel struct {
	ApiManagementId  string `tfschema:"api_management_id"`
	NotificationName string `tfschema:"notification_type"`
	UserId           string `tfschema:"user_id"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementNotificationRecipientUserModel{}
}
