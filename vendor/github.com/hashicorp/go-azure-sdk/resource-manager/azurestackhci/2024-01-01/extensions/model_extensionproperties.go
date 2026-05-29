package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExtensionProperties struct {
	AggregateState          *ExtensionAggregateState `json:"aggregateState,omitempty"`
	ExtensionParameters     *ExtensionParameters     `json:"extensionParameters,omitempty"`
	ManagedBy               *ExtensionManagedBy      `json:"managedBy,omitempty"`
	PerNodeExtensionDetails *[]PerNodeExtensionState `json:"perNodeExtensionDetails,omitempty"`
	ProvisioningState       *ProvisioningState       `json:"provisioningState,omitempty"`
}
