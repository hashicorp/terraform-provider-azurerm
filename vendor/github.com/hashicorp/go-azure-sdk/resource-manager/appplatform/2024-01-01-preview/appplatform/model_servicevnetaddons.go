package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceVNetAddons struct {
	DataPlanePublicEndpoint *bool `json:"dataPlanePublicEndpoint,omitempty"`
	LogStreamPublicEndpoint *bool `json:"logStreamPublicEndpoint,omitempty"`
}
