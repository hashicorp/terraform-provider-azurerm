package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopologyResource struct {
	Associations *[]TopologyAssociation `json:"associations,omitempty"`
	Id           *string                `json:"id,omitempty"`
	Location     *string                `json:"location,omitempty"`
	Name         *string                `json:"name,omitempty"`
}
