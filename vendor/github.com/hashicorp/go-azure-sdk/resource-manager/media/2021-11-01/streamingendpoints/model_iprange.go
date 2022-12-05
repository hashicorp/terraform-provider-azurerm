package streamingendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPRange struct {
	Address            *string `json:"address,omitempty"`
	Name               *string `json:"name,omitempty"`
	SubnetPrefixLength *int64  `json:"subnetPrefixLength,omitempty"`
}
