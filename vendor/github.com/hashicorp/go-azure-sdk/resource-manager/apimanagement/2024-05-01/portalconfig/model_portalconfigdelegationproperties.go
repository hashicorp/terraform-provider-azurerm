package portalconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalConfigDelegationProperties struct {
	DelegateRegistration *bool   `json:"delegateRegistration,omitempty"`
	DelegateSubscription *bool   `json:"delegateSubscription,omitempty"`
	DelegationURL        *string `json:"delegationUrl,omitempty"`
	ValidationKey        *string `json:"validationKey,omitempty"`
}
