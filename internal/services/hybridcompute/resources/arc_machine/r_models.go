// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ArcMachineResourceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	Kind              string                         `tfschema:"kind"`
	Identity          []identity.ModelSystemAssigned `tfschema:"identity"`
	Tags              map[string]string              `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &ArcMachineResourceModel{}
}
