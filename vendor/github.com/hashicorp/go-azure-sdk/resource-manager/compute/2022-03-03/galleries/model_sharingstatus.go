package galleries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingStatus struct {
	AggregatedState *SharingState            `json:"aggregatedState,omitempty"`
	Summary         *[]RegionalSharingStatus `json:"summary,omitempty"`
}
