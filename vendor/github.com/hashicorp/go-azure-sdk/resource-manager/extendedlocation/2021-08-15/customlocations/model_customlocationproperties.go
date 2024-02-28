package customlocations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomLocationProperties struct {
	Authentication      *CustomLocationPropertiesAuthentication `json:"authentication,omitempty"`
	ClusterExtensionIds *[]string                               `json:"clusterExtensionIds,omitempty"`
	DisplayName         *string                                 `json:"displayName,omitempty"`
	HostResourceId      *string                                 `json:"hostResourceId,omitempty"`
	HostType            *HostType                               `json:"hostType,omitempty"`
	Namespace           *string                                 `json:"namespace,omitempty"`
	ProvisioningState   *string                                 `json:"provisioningState,omitempty"`
}
