package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointDestination struct {
	ServiceResourceId *string     `json:"serviceResourceId,omitempty"`
	SparkEnabled      *bool       `json:"sparkEnabled,omitempty"`
	SparkStatus       *RuleStatus `json:"sparkStatus,omitempty"`
	SubresourceTarget *string     `json:"subresourceTarget,omitempty"`
}
