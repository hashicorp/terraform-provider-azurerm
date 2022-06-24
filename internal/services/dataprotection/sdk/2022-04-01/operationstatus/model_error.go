package operationstatus

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Error struct {
	AdditionalInfo *[]ErrorAdditionalInfo `json:"additionalInfo,omitempty"`
	Code           *string                `json:"code,omitempty"`
	Details        *[]Error               `json:"details,omitempty"`
	Message        *string                `json:"message,omitempty"`
	Target         *string                `json:"target,omitempty"`
}
