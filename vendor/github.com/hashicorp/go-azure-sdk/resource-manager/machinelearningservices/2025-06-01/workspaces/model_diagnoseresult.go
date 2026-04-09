package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiagnoseResult struct {
	Code    *string              `json:"code,omitempty"`
	Level   *DiagnoseResultLevel `json:"level,omitempty"`
	Message *string              `json:"message,omitempty"`
}
