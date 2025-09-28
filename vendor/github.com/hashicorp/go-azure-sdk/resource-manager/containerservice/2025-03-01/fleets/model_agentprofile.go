package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentProfile struct {
	SubnetId *string `json:"subnetId,omitempty"`
	VMSize   *string `json:"vmSize,omitempty"`
}
