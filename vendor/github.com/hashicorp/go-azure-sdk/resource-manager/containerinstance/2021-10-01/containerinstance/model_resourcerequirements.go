package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceRequirements struct {
	Limits   *ResourceLimits  `json:"limits,omitempty"`
	Requests ResourceRequests `json:"requests"`
}
