package imagedefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDefinitionProperties struct {
	ActiveImageReference        *ImageReference                  `json:"activeImageReference,omitempty"`
	AutoImageBuild              *AutoImageBuildStatus            `json:"autoImageBuild,omitempty"`
	FileURL                     *string                          `json:"fileUrl,omitempty"`
	ImageReference              *ImageReference                  `json:"imageReference,omitempty"`
	ImageValidationErrorDetails *ImageValidationErrorDetails     `json:"imageValidationErrorDetails,omitempty"`
	ImageValidationStatus       *ImageValidationStatus           `json:"imageValidationStatus,omitempty"`
	LatestBuild                 *LatestImageBuild                `json:"latestBuild,omitempty"`
	ValidationStatus            *CatalogResourceValidationStatus `json:"validationStatus,omitempty"`
}
