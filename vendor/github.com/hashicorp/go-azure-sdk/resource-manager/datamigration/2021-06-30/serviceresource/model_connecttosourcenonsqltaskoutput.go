package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToSourceNonSqlTaskOutput struct {
	Databases                *[]string              `json:"databases,omitempty"`
	Id                       *string                `json:"id,omitempty"`
	ServerProperties         *ServerProperties      `json:"serverProperties,omitempty"`
	SourceServerBrandVersion *string                `json:"sourceServerBrandVersion,omitempty"`
	ValidationErrors         *[]ReportableException `json:"validationErrors,omitempty"`
}
