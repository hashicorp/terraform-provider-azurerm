package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchingConfiguration struct {
	LatencySeconds *int64 `json:"latencySeconds,omitempty"`
	MaxMessages    *int64 `json:"maxMessages,omitempty"`
}
