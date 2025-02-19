package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HopLink struct {
	Context    *map[string]string   `json:"context,omitempty"`
	Issues     *[]ConnectivityIssue `json:"issues,omitempty"`
	LinkType   *string              `json:"linkType,omitempty"`
	NextHopId  *string              `json:"nextHopId,omitempty"`
	Properties *HopLinkProperties   `json:"properties,omitempty"`
	ResourceId *string              `json:"resourceId,omitempty"`
}
