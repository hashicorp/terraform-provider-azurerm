package containerappssessionpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionPoolProperties struct {
	ContainerType               *ContainerType                `json:"containerType,omitempty"`
	CustomContainerTemplate     *CustomContainerTemplate      `json:"customContainerTemplate,omitempty"`
	DynamicPoolConfiguration    *DynamicPoolConfiguration     `json:"dynamicPoolConfiguration,omitempty"`
	EnvironmentId               *string                       `json:"environmentId,omitempty"`
	ManagedIdentitySettings     *[]ManagedIdentitySetting     `json:"managedIdentitySettings,omitempty"`
	NodeCount                   *int64                        `json:"nodeCount,omitempty"`
	PoolManagementEndpoint      *string                       `json:"poolManagementEndpoint,omitempty"`
	PoolManagementType          *PoolManagementType           `json:"poolManagementType,omitempty"`
	ProvisioningState           *SessionPoolProvisioningState `json:"provisioningState,omitempty"`
	ScaleConfiguration          *ScaleConfiguration           `json:"scaleConfiguration,omitempty"`
	Secrets                     *[]SessionPoolSecret          `json:"secrets,omitempty"`
	SessionNetworkConfiguration *SessionNetworkConfiguration  `json:"sessionNetworkConfiguration,omitempty"`
}
