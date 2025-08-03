package containerappsrevisions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Scale struct {
	CooldownPeriod  *int64       `json:"cooldownPeriod,omitempty"`
	MaxReplicas     *int64       `json:"maxReplicas,omitempty"`
	MinReplicas     *int64       `json:"minReplicas,omitempty"`
	PollingInterval *int64       `json:"pollingInterval,omitempty"`
	Rules           *[]ScaleRule `json:"rules,omitempty"`
}
