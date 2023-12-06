package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopologyAssociation struct {
	AssociationType *AssociationType `json:"associationType,omitempty"`
	Name            *string          `json:"name,omitempty"`
	ResourceId      *string          `json:"resourceId,omitempty"`
}
