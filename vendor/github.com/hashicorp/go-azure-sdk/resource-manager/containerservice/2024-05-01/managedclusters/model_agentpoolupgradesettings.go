package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolUpgradeSettings struct {
	DrainTimeoutInMinutes     *int64  `json:"drainTimeoutInMinutes,omitempty"`
	MaxSurge                  *string `json:"maxSurge,omitempty"`
	NodeSoakDurationInMinutes *int64  `json:"nodeSoakDurationInMinutes,omitempty"`
}
