package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverConnectionDetails struct {
	FailoverConnectionName *string `json:"failoverConnectionName,omitempty"`
	FailoverLocation       *string `json:"failoverLocation,omitempty"`
	IsVerified             *bool   `json:"isVerified,omitempty"`
}
