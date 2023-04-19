package simpolicy

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SimPolicyPropertiesFormat struct {
	DefaultSlice          SliceResourceId                   `json:"defaultSlice"`
	ProvisioningState     *ProvisioningState                `json:"provisioningState,omitempty"`
	RegistrationTimer     *int64                            `json:"registrationTimer,omitempty"`
	RfspIndex             *int64                            `json:"rfspIndex,omitempty"`
	SiteProvisioningState *map[string]SiteProvisioningState `json:"siteProvisioningState,omitempty"`
	SliceConfigurations   []SliceConfiguration              `json:"sliceConfigurations"`
	UeAmbr                Ambr                              `json:"ueAmbr"`
}
