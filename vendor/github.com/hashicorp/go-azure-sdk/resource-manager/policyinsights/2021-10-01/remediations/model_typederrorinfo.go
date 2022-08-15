package remediations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TypedErrorInfo struct {
	Info *interface{} `json:"info,omitempty"`
	Type *string      `json:"type,omitempty"`
}
