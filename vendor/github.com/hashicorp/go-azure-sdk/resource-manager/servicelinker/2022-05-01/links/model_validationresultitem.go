package links

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationResultItem struct {
	Description  *string                 `json:"description,omitempty"`
	ErrorCode    *string                 `json:"errorCode,omitempty"`
	ErrorMessage *string                 `json:"errorMessage,omitempty"`
	Name         *string                 `json:"name,omitempty"`
	Result       *ValidationResultStatus `json:"result,omitempty"`
}
