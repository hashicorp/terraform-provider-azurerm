// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
)

type ApplianceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	Distro            appliances.Distro              `tfschema:"distro"`
	Identity          []identity.ModelSystemAssigned `tfschema:"identity"`
	Provider          appliances.Provider            `tfschema:"infrastructure_provider"`
	PublicKeyBase64   string                         `tfschema:"public_key_base64"`
	Tags              map[string]any                 `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &Resource{}
}
