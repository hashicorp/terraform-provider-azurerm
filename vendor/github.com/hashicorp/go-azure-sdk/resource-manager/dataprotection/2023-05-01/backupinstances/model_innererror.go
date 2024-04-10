package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InnerError struct {
	AdditionalInfo     *map[string]string `json:"additionalInfo,omitempty"`
	Code               *string            `json:"code,omitempty"`
	EmbeddedInnerError *InnerError        `json:"embeddedInnerError,omitempty"`
}
