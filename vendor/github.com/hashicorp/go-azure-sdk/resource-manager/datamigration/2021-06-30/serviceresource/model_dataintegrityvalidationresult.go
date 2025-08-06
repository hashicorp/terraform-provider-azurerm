package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataIntegrityValidationResult struct {
	FailedObjects    *map[string]string `json:"failedObjects,omitempty"`
	ValidationErrors *ValidationError   `json:"validationErrors,omitempty"`
}
