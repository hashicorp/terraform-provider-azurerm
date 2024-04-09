package sapcentralinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnqueueServerProperties struct {
	Health    *SAPHealthState `json:"health,omitempty"`
	Hostname  *string         `json:"hostname,omitempty"`
	IPAddress *string         `json:"ipAddress,omitempty"`
	Port      *int64          `json:"port,omitempty"`
}
