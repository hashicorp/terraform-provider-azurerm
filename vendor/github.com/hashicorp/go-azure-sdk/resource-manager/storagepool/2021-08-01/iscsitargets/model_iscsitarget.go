package iscsitargets

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IscsiTarget struct {
	Id                *string                `json:"id,omitempty"`
	ManagedBy         *string                `json:"managedBy,omitempty"`
	ManagedByExtended *[]string              `json:"managedByExtended,omitempty"`
	Name              *string                `json:"name,omitempty"`
	Properties        IscsiTargetProperties  `json:"properties"`
	SystemData        *systemdata.SystemData `json:"systemData,omitempty"`
	Type              *string                `json:"type,omitempty"`
}
