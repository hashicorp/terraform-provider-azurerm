package hybridrunbookworkergroup

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridRunbookWorkerGroup struct {
	Credential           *RunAsCredentialAssociationProperty `json:"credential,omitempty"`
	GroupType            *GroupTypeEnum                      `json:"groupType,omitempty"`
	HybridRunbookWorkers *[]HybridRunbookWorkerLegacy        `json:"hybridRunbookWorkers,omitempty"`
	Id                   *string                             `json:"id,omitempty"`
	Name                 *string                             `json:"name,omitempty"`
	SystemData           *systemdata.SystemData              `json:"systemData,omitempty"`
	Type                 *string                             `json:"type,omitempty"`
}
