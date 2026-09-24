package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosCustomPolicyPropertiesFormat struct {
	DetectionRules          *[]DdosDetectionRule `json:"detectionRules,omitempty"`
	FrontEndIPConfiguration *[]SubResource       `json:"frontEndIpConfiguration,omitempty"`
	ProvisioningState       *ProvisioningState   `json:"provisioningState,omitempty"`
	PublicIPAddresses       *[]SubResource       `json:"publicIPAddresses,omitempty"`
	ResourceGuid            *string              `json:"resourceGuid,omitempty"`
}
