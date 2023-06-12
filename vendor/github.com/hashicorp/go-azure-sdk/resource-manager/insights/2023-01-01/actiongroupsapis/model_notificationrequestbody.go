package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationRequestBody struct {
	AlertType                  string                       `json:"alertType"`
	ArmRoleReceivers           *[]ArmRoleReceiver           `json:"armRoleReceivers,omitempty"`
	AutomationRunbookReceivers *[]AutomationRunbookReceiver `json:"automationRunbookReceivers,omitempty"`
	AzureAppPushReceivers      *[]AzureAppPushReceiver      `json:"azureAppPushReceivers,omitempty"`
	AzureFunctionReceivers     *[]AzureFunctionReceiver     `json:"azureFunctionReceivers,omitempty"`
	EmailReceivers             *[]EmailReceiver             `json:"emailReceivers,omitempty"`
	EventHubReceivers          *[]EventHubReceiver          `json:"eventHubReceivers,omitempty"`
	ItsmReceivers              *[]ItsmReceiver              `json:"itsmReceivers,omitempty"`
	LogicAppReceivers          *[]LogicAppReceiver          `json:"logicAppReceivers,omitempty"`
	SmsReceivers               *[]SmsReceiver               `json:"smsReceivers,omitempty"`
	VoiceReceivers             *[]VoiceReceiver             `json:"voiceReceivers,omitempty"`
	WebhookReceivers           *[]WebhookReceiver           `json:"webhookReceivers,omitempty"`
}
