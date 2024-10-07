package verifiedpartners

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VerifiedPartnerProperties struct {
	OrganizationName               *string                           `json:"organizationName,omitempty"`
	PartnerDisplayName             *string                           `json:"partnerDisplayName,omitempty"`
	PartnerRegistrationImmutableId *string                           `json:"partnerRegistrationImmutableId,omitempty"`
	PartnerTopicDetails            *PartnerDetails                   `json:"partnerTopicDetails,omitempty"`
	ProvisioningState              *VerifiedPartnerProvisioningState `json:"provisioningState,omitempty"`
}
