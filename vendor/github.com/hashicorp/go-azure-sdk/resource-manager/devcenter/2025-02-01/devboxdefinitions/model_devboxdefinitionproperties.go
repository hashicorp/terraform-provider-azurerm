package devboxdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DevBoxDefinitionProperties struct {
	ActiveImageReference        *ImageReference                  `json:"activeImageReference,omitempty"`
	HibernateSupport            *HibernateSupport                `json:"hibernateSupport,omitempty"`
	ImageReference              *ImageReference                  `json:"imageReference,omitempty"`
	ImageValidationErrorDetails *ImageValidationErrorDetails     `json:"imageValidationErrorDetails,omitempty"`
	ImageValidationStatus       *ImageValidationStatus           `json:"imageValidationStatus,omitempty"`
	OsStorageType               *string                          `json:"osStorageType,omitempty"`
	ProvisioningState           *ProvisioningState               `json:"provisioningState,omitempty"`
	Sku                         *Sku                             `json:"sku,omitempty"`
	ValidationStatus            *CatalogResourceValidationStatus `json:"validationStatus,omitempty"`
}
