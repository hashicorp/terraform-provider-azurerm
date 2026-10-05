package resourceanchors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceAnchorProperties struct {
	LinkedCompartmentId *string                         `json:"linkedCompartmentId,omitempty"`
	ProvisioningState   *AzureResourceProvisioningState `json:"provisioningState,omitempty"`
}
