package containerappssessionpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleConfiguration struct {
	MaxConcurrentSessions *int64 `json:"maxConcurrentSessions,omitempty"`
	ReadySessionInstances *int64 `json:"readySessionInstances,omitempty"`
}
