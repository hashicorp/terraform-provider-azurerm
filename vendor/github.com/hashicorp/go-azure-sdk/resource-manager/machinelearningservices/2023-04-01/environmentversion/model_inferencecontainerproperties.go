package environmentversion

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InferenceContainerProperties struct {
	LivenessRoute  *Route `json:"livenessRoute,omitempty"`
	ReadinessRoute *Route `json:"readinessRoute,omitempty"`
	ScoringRoute   *Route `json:"scoringRoute,omitempty"`
}
