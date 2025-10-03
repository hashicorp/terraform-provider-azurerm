package partnerconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerConfigurationProperties struct {
	PartnerAuthorization *PartnerAuthorization                  `json:"partnerAuthorization,omitempty"`
	ProvisioningState    *PartnerConfigurationProvisioningState `json:"provisioningState,omitempty"`
}
