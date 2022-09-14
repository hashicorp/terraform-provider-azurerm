package privatezones

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateZoneProperties struct {
	MaxNumberOfRecordSets                          *int64             `json:"maxNumberOfRecordSets,omitempty"`
	MaxNumberOfVirtualNetworkLinks                 *int64             `json:"maxNumberOfVirtualNetworkLinks,omitempty"`
	MaxNumberOfVirtualNetworkLinksWithRegistration *int64             `json:"maxNumberOfVirtualNetworkLinksWithRegistration,omitempty"`
	NumberOfRecordSets                             *int64             `json:"numberOfRecordSets,omitempty"`
	NumberOfVirtualNetworkLinks                    *int64             `json:"numberOfVirtualNetworkLinks,omitempty"`
	NumberOfVirtualNetworkLinksWithRegistration    *int64             `json:"numberOfVirtualNetworkLinksWithRegistration,omitempty"`
	ProvisioningState                              *ProvisioningState `json:"provisioningState,omitempty"`
}
