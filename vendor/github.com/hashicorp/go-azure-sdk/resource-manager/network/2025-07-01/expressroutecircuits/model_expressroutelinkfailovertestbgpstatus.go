package expressroutecircuits

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteLinkFailoverTestBgpStatus struct {
	CheckTimeUtc *string                                     `json:"checkTimeUtc,omitempty"`
	Link         *ExpressRouteFailoverLinkType               `json:"link,omitempty"`
	Status       *ExpressRouteLinkFailoverBgpStatus          `json:"status,omitempty"`
	Type         *ExpressRouteFailoverBgpStatusAddressFamily `json:"type,omitempty"`
}
