package resourcegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceGroupExportResult struct {
	Error    *ErrorResponse `json:"error,omitempty"`
	Template *interface{}   `json:"template,omitempty"`
}
