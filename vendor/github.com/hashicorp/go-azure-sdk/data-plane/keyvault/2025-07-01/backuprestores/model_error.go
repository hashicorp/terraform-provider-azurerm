package backuprestores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Error struct {
	Code       *string `json:"code,omitempty"`
	Innererror *Error  `json:"innererror,omitempty"`
	Message    *string `json:"message,omitempty"`
}
