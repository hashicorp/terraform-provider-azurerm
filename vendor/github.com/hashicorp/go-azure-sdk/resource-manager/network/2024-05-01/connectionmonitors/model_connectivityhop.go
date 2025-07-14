package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityHop struct {
	Address        *string              `json:"address,omitempty"`
	Id             *string              `json:"id,omitempty"`
	Issues         *[]ConnectivityIssue `json:"issues,omitempty"`
	Links          *[]HopLink           `json:"links,omitempty"`
	NextHopIds     *[]string            `json:"nextHopIds,omitempty"`
	PreviousHopIds *[]string            `json:"previousHopIds,omitempty"`
	PreviousLinks  *[]HopLink           `json:"previousLinks,omitempty"`
	ResourceId     *string              `json:"resourceId,omitempty"`
	Type           *string              `json:"type,omitempty"`
}
