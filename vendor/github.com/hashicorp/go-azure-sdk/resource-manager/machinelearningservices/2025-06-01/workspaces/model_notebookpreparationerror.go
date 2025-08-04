package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotebookPreparationError struct {
	ErrorMessage *string `json:"errorMessage,omitempty"`
	StatusCode   *int64  `json:"statusCode,omitempty"`
}
