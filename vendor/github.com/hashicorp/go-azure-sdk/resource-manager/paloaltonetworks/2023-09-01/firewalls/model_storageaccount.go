package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageAccount struct {
	AccountName    *string `json:"accountName,omitempty"`
	Id             *string `json:"id,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}
