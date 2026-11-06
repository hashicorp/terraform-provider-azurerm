package invoicesection

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InvoiceSectionProperties struct {
	DisplayName       *string                        `json:"displayName,omitempty"`
	ProvisioningState *ProvisioningState             `json:"provisioningState,omitempty"`
	ReasonCode        *InvoiceSectionStateReasonCode `json:"reasonCode,omitempty"`
	State             *InvoiceSectionState           `json:"state,omitempty"`
	SystemId          *string                        `json:"systemId,omitempty"`
	Tags              *map[string]string             `json:"tags,omitempty"`
	TargetCloud       *string                        `json:"targetCloud,omitempty"`
}
