package signalr

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpstreamTemplate struct {
	Auth            *UpstreamAuthSettings `json:"auth,omitempty"`
	CategoryPattern *string               `json:"categoryPattern,omitempty"`
	EventPattern    *string               `json:"eventPattern,omitempty"`
	HubPattern      *string               `json:"hubPattern,omitempty"`
	UrlTemplate     string                `json:"urlTemplate"`
}
