package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationStatus struct {
	AggregatedState *AggregatedReplicationState  `json:"aggregatedState,omitempty"`
	Summary         *[]RegionalReplicationStatus `json:"summary,omitempty"`
}
