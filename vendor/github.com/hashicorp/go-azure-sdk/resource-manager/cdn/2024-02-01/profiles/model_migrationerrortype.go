package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationErrorType struct {
	Code         *string `json:"code,omitempty"`
	ErrorMessage *string `json:"errorMessage,omitempty"`
	NextSteps    *string `json:"nextSteps,omitempty"`
	ResourceName *string `json:"resourceName,omitempty"`
}
