package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventPreview struct {
	AccessControl       *LiveEventPreviewAccessControl `json:"accessControl,omitempty"`
	AlternativeMediaId  *string                        `json:"alternativeMediaId,omitempty"`
	Endpoints           *[]LiveEventEndpoint           `json:"endpoints,omitempty"`
	PreviewLocator      *string                        `json:"previewLocator,omitempty"`
	StreamingPolicyName *string                        `json:"streamingPolicyName,omitempty"`
}
