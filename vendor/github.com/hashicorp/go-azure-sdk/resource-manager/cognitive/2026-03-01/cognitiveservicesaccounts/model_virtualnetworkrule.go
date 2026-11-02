package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkRule struct {
	Id                               string  `json:"id"`
	IgnoreMissingVnetServiceEndpoint *bool   `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
	State                            *string `json:"state,omitempty"`
}
