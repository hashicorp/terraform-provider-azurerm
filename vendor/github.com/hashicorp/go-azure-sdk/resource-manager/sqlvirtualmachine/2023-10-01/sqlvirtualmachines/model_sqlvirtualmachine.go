package sqlvirtualmachines

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SqlVirtualMachine struct {
	Id         *string                      `json:"id,omitempty"`
	Identity   *ResourceIdentity            `json:"identity,omitempty"`
	Location   string                       `json:"location"`
	Name       *string                      `json:"name,omitempty"`
	Properties *SqlVirtualMachineProperties `json:"properties,omitempty"`
	SystemData *systemdata.SystemData       `json:"systemData,omitempty"`
	Tags       *map[string]string           `json:"tags,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}
