package firstpartyservicetags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirstPartyServiceTagPropertiesFormat struct {
	FailedReason      *string            `json:"failedReason,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceGuid      *string            `json:"resourceGuid,omitempty"`
	Value             string             `json:"value"`
}
