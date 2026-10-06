package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceOracleSyncTaskOutput struct {
	Databases                *[]string              `json:"databases,omitempty"`
	SourceServerBrandVersion *string                `json:"sourceServerBrandVersion,omitempty"`
	SourceServerVersion      *string                `json:"sourceServerVersion,omitempty"`
	ValidationErrors         *[]ReportableException `json:"validationErrors,omitempty"`
}
