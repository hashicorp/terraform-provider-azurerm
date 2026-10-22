// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_exadata_infrastructure

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

type ExadataInfraResourceModel struct {
	Location          string            `tfschema:"location"`
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Tags              map[string]string `tfschema:"tags"`

	ComputeCount int64  `tfschema:"compute_count"`
	DisplayName  string `tfschema:"display_name"`
	Shape        string `tfschema:"shape"`
	StorageCount int64  `tfschema:"storage_count"`

	CustomerContacts   []string                 `tfschema:"customer_contacts"`
	DatabaseServerType string                   `tfschema:"database_server_type"`
	MaintenanceWindow  []MaintenanceWindowModel `tfschema:"maintenance_window"`
	StorageServerType  string                   `tfschema:"storage_server_type"`
	Zones              zones.Schema             `tfschema:"zones"`
}

func (Resource) ModelObject() any {
	return &Resource{}
}
