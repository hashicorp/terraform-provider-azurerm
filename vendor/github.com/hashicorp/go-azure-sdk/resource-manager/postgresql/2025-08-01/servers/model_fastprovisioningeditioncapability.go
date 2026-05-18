package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FastProvisioningEditionCapability struct {
	Reason                  *string           `json:"reason,omitempty"`
	ServerCount             *int64            `json:"serverCount,omitempty"`
	Status                  *CapabilityStatus `json:"status,omitempty"`
	SupportedServerVersions *string           `json:"supportedServerVersions,omitempty"`
	SupportedSku            *string           `json:"supportedSku,omitempty"`
	SupportedStorageGb      *int64            `json:"supportedStorageGb,omitempty"`
	SupportedTier           *string           `json:"supportedTier,omitempty"`
}
