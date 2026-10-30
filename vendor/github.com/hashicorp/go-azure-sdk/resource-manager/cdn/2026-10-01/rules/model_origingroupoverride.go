package rules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OriginGroupOverride struct {
	ForwardingProtocol *ForwardingProtocol `json:"forwardingProtocol,omitempty"`
	OriginGroup        *ResourceReference  `json:"originGroup,omitempty"`
}
