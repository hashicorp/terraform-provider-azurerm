package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CalloutPolicy struct {
	CalloutId       *string         `json:"calloutId,omitempty"`
	CalloutType     *CalloutType    `json:"calloutType,omitempty"`
	CalloutUriRegex *string         `json:"calloutUriRegex,omitempty"`
	OutboundAccess  *OutboundAccess `json:"outboundAccess,omitempty"`
}
