// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type ServerModel struct {
	Name               string                                     `tfschema:"name"`
	ResourceGroup      string                                     `tfschema:"resource_group_name"`
	Location           string                                     `tfschema:"location"`
	StorageSKU         string                                     `tfschema:"storage_sku"`
	FrsTenantId        string                                     `tfschema:"frs_tenant_id"`
	PrimaryKey         string                                     `tfschema:"primary_key"`
	SecondaryKey       string                                     `tfschema:"secondary_key"`
	OrdererEndpoints   []string                                   `tfschema:"orderer_endpoints"`
	StorageEndpoints   []string                                   `tfschema:"storage_endpoints"`
	ServiceEndpoints   []string                                   `tfschema:"service_endpoints"`
	Tags               map[string]string                          `tfschema:"tags"`
	Identity           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	CustomerManagedKey []CustomerManagedKey                       `tfschema:"customer_managed_key"`
}

func (s Resource) ModelObject() any {
	return &ServerModel{}
}
