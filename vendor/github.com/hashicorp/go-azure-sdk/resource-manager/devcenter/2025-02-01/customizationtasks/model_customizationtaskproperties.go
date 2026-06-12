package customizationtasks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomizationTaskProperties struct {
	Inputs           *map[string]CustomizationTaskInput `json:"inputs,omitempty"`
	Timeout          *int64                             `json:"timeout,omitempty"`
	ValidationStatus *CatalogResourceValidationStatus   `json:"validationStatus,omitempty"`
}
