package connectedregistries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedRegistryUpdateProperties struct {
	ClientTokenIds    *[]string             `json:"clientTokenIds,omitempty"`
	Logging           *LoggingProperties    `json:"logging,omitempty"`
	NotificationsList *[]string             `json:"notificationsList,omitempty"`
	SyncProperties    *SyncUpdateProperties `json:"syncProperties,omitempty"`
}
