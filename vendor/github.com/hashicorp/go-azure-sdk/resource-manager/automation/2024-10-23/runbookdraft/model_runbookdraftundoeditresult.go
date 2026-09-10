package runbookdraft

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookDraftUndoEditResult struct {
	RequestId  *string         `json:"requestId,omitempty"`
	StatusCode *HTTPStatusCode `json:"statusCode,omitempty"`
}
