package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResizeError struct {
	Code    string         `json:"code"`
	Details *[]ResizeError `json:"details,omitempty"`
	Message string         `json:"message"`
}
