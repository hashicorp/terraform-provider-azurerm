package managedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedHSMSecurityDomainProperties struct {
	ActivationStatus        *ActivationStatus `json:"activationStatus,omitempty"`
	ActivationStatusMessage *string           `json:"activationStatusMessage,omitempty"`
}
