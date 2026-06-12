package galleryimageversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyViolation struct {
	Category *PolicyViolationCategory `json:"category,omitempty"`
	Details  *string                  `json:"details,omitempty"`
}
