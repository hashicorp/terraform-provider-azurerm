package protectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenericContainerExtendedInfo struct {
	ContainerIdentityInfo *ContainerIdentityInfo `json:"containerIdentityInfo,omitempty"`
	RawCertData           *string                `json:"rawCertData,omitempty"`
	ServiceEndpoints      *map[string]string     `json:"serviceEndpoints,omitempty"`
}
