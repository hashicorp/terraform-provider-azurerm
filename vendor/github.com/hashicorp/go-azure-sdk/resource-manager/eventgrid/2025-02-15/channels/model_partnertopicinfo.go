package channels

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerTopicInfo struct {
	AzureSubscriptionId *string        `json:"azureSubscriptionId,omitempty"`
	EventTypeInfo       *EventTypeInfo `json:"eventTypeInfo,omitempty"`
	Name                *string        `json:"name,omitempty"`
	ResourceGroupName   *string        `json:"resourceGroupName,omitempty"`
	Source              *string        `json:"source,omitempty"`
}
