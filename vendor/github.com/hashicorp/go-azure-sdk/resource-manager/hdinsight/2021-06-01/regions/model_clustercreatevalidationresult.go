package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterCreateValidationResult struct {
	AaddsResourcesDetails     *[]AaddsResourceDetails `json:"aaddsResourcesDetails,omitempty"`
	EstimatedCreationDuration *string                 `json:"estimatedCreationDuration,omitempty"`
	ValidationErrors          *[]ValidationErrorInfo  `json:"validationErrors,omitempty"`
	ValidationWarnings        *[]ValidationErrorInfo  `json:"validationWarnings,omitempty"`
}
