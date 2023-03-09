package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AsyncOperationResult struct {
	Error  *ErrorDetail          `json:"error,omitempty"`
	Name   *string               `json:"name,omitempty"`
	Status *AsyncOperationStatus `json:"status,omitempty"`
}
