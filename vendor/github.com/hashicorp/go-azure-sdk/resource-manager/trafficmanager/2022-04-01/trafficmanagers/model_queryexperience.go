package trafficmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryExperience struct {
	EndpointId int64    `json:"endpointId"`
	Latency    *float64 `json:"latency,omitempty"`
	QueryCount int64    `json:"queryCount"`
}
