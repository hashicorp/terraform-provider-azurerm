package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConflictResolutionPolicy struct {
	ConflictResolutionPath      *string                 `json:"conflictResolutionPath,omitempty"`
	ConflictResolutionProcedure *string                 `json:"conflictResolutionProcedure,omitempty"`
	Mode                        *ConflictResolutionMode `json:"mode,omitempty"`
}
