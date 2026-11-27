// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_discovery_virtual_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type WorkloadsSAPDiscoveryVirtualInstanceModel struct {
	Name                              string                       `tfschema:"name"`
	ResourceGroupName                 string                       `tfschema:"resource_group_name"`
	Location                          string                       `tfschema:"location"`
	CentralServerVmId                 string                       `tfschema:"central_server_virtual_machine_id"`
	Environment                       string                       `tfschema:"environment"`
	Identity                          []identity.ModelUserAssigned `tfschema:"identity"`
	ManagedResourceGroupName          string                       `tfschema:"managed_resource_group_name"`
	ManagedResourcesNetworkAccessType string                       `tfschema:"managed_resources_network_access_type"`
	ManagedStorageAccountName         string                       `tfschema:"managed_storage_account_name"`
	SapProduct                        string                       `tfschema:"sap_product"`
	Tags                              map[string]string            `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &WorkloadsSAPDiscoveryVirtualInstanceModel{}
}
