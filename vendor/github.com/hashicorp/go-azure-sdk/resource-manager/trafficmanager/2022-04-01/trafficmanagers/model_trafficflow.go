package trafficmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficFlow struct {
	Latitude         *float64           `json:"latitude,omitempty"`
	Longitude        *float64           `json:"longitude,omitempty"`
	QueryExperiences *[]QueryExperience `json:"queryExperiences,omitempty"`
	SourceIP         *string            `json:"sourceIp,omitempty"`
}
