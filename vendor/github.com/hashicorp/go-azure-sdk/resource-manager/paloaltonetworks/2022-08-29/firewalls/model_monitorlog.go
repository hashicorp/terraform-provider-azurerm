package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorLog struct {
	Id             *string `json:"id,omitempty"`
	PrimaryKey     *string `json:"primaryKey,omitempty"`
	SecondaryKey   *string `json:"secondaryKey,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
	Workspace      *string `json:"workspace,omitempty"`
}
