package partnerregistrations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerRegistrationProperties struct {
	PartnerRegistrationImmutableId *string                               `json:"partnerRegistrationImmutableId,omitempty"`
	ProvisioningState              *PartnerRegistrationProvisioningState `json:"provisioningState,omitempty"`
}
