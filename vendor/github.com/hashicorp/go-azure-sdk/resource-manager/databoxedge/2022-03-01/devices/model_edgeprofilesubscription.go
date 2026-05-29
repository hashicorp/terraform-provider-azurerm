package devices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeProfileSubscription struct {
	Id               *string                 `json:"id,omitempty"`
	Properties       *SubscriptionProperties `json:"properties,omitempty"`
	RegistrationDate *string                 `json:"registrationDate,omitempty"`
	RegistrationId   *string                 `json:"registrationId,omitempty"`
	State            *SubscriptionState      `json:"state,omitempty"`
	SubscriptionId   *string                 `json:"subscriptionId,omitempty"`
}
