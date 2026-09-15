package containerappssessionpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionPoolUpdatablePropertiesProperties struct {
	CustomContainerTemplate     *CustomContainerTemplate     `json:"customContainerTemplate,omitempty"`
	DynamicPoolConfiguration    *DynamicPoolConfiguration    `json:"dynamicPoolConfiguration,omitempty"`
	ScaleConfiguration          *ScaleConfiguration          `json:"scaleConfiguration,omitempty"`
	Secrets                     *[]SessionPoolSecret         `json:"secrets,omitempty"`
	SessionNetworkConfiguration *SessionNetworkConfiguration `json:"sessionNetworkConfiguration,omitempty"`
}
