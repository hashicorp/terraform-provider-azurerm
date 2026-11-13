package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PeeringHealth struct {
	StatsAtPrepare *PeeringStats `json:"statsAtPrepare,omitempty"`
	StatsCurrent   *PeeringStats `json:"statsCurrent,omitempty"`
	Type           *string       `json:"type,omitempty"`
}
