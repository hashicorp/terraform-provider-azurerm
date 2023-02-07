package snapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyCompletionError struct {
	ErrorCode    CopyCompletionErrorReason `json:"errorCode"`
	ErrorMessage string                    `json:"errorMessage"`
}
