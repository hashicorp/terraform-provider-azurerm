package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceMoveDetails struct {
	CompletionTimeUtc  *string `json:"completionTimeUtc,omitempty"`
	OperationId        *string `json:"operationId,omitempty"`
	SourceResourcePath *string `json:"sourceResourcePath,omitempty"`
	StartTimeUtc       *string `json:"startTimeUtc,omitempty"`
	TargetResourcePath *string `json:"targetResourcePath,omitempty"`
}
