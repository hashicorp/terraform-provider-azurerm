package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetUserTablesMySqlTaskOutput struct {
	DatabasesToTables *map[string][]DatabaseTable `json:"databasesToTables,omitempty"`
	Id                *string                     `json:"id,omitempty"`
	ValidationErrors  *[]ReportableException      `json:"validationErrors,omitempty"`
}
