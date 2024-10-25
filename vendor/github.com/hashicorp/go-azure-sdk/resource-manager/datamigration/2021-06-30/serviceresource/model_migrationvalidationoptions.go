package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrationValidationOptions struct {
	EnableDataIntegrityValidation *bool `json:"enableDataIntegrityValidation,omitempty"`
	EnableQueryAnalysisValidation *bool `json:"enableQueryAnalysisValidation,omitempty"`
	EnableSchemaValidation        *bool `json:"enableSchemaValidation,omitempty"`
}
