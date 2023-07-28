package batchdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BatchRetrySettings struct {
	MaxRetries *int64  `json:"maxRetries,omitempty"`
	Timeout    *string `json:"timeout,omitempty"`
}
