package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHub struct {
	Id             *string `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	NameSpace      *string `json:"nameSpace,omitempty"`
	PolicyName     *string `json:"policyName,omitempty"`
	SubscriptionId *string `json:"subscriptionId,omitempty"`
}
