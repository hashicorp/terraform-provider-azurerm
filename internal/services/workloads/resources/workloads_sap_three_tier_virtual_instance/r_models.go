// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type WorkloadsSAPThreeTierVirtualInstanceModel struct {
	Name                              string                       `tfschema:"name"`
	ResourceGroupName                 string                       `tfschema:"resource_group_name"`
	Location                          string                       `tfschema:"location"`
	AppLocation                       string                       `tfschema:"app_location"`
	Environment                       string                       `tfschema:"environment"`
	Identity                          []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName          string                       `tfschema:"managed_resource_group_name"`
	ManagedResourcesNetworkAccessType string                       `tfschema:"managed_resources_network_access_type"`
	SapFqdn                           string                       `tfschema:"sap_fqdn"`
	SapProduct                        string                       `tfschema:"sap_product"`
	ThreeTierConfiguration            []ThreeTierConfiguration     `tfschema:"three_tier_configuration"`
	Tags                              map[string]string            `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &WorkloadsSAPThreeTierVirtualInstanceModel{}
}
