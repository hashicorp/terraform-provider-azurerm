// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_source_endpoint

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/endpoints"
)

type StorageMoverSourceEndpointModel struct {
	Name           string               `tfschema:"name"`
	StorageMoverId string               `tfschema:"storage_mover_id"`
	Export         string               `tfschema:"export"`
	Host           string               `tfschema:"host"`
	NfsVersion     endpoints.NfsVersion `tfschema:"nfs_version"`
	Description    string               `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &StorageMoverSourceEndpointModel{}
}
