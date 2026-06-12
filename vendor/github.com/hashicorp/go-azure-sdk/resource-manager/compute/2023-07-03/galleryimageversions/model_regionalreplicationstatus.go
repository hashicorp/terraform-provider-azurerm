package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegionalReplicationStatus struct {
	Details  *string           `json:"details,omitempty"`
	Progress *int64            `json:"progress,omitempty"`
	Region   *string           `json:"region,omitempty"`
	State    *ReplicationState `json:"state,omitempty"`
}
