package storagediscoveryworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportResultSet struct {
	Columns   *[]ReportResultColumn `json:"columns,omitempty"`
	ErrorCode *string               `json:"errorCode,omitempty"`
	Rows      *[][]string           `json:"rows,omitempty"`
}
