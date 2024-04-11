package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayApiMetadataProperties struct {
	Description   *string `json:"description,omitempty"`
	Documentation *string `json:"documentation,omitempty"`
	ServerUrl     *string `json:"serverUrl,omitempty"`
	Title         *string `json:"title,omitempty"`
	Version       *string `json:"version,omitempty"`
}
