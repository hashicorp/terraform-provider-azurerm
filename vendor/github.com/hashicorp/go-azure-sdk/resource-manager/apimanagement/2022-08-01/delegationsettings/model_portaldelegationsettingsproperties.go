package delegationsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalDelegationSettingsProperties struct {
	Subscriptions    *SubscriptionsDelegationSettingsProperties `json:"subscriptions,omitempty"`
	Url              *string                                    `json:"url,omitempty"`
	UserRegistration *RegistrationDelegationSettingsProperties  `json:"userRegistration,omitempty"`
	ValidationKey    *string                                    `json:"validationKey,omitempty"`
}
