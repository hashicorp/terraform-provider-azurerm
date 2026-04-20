package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryImageVersionSafetyProfile struct {
	AllowDeletionOfReplicatedLocations *bool              `json:"allowDeletionOfReplicatedLocations,omitempty"`
	PolicyViolations                   *[]PolicyViolation `json:"policyViolations,omitempty"`
	ReportedForPolicyViolation         *bool              `json:"reportedForPolicyViolation,omitempty"`
}
