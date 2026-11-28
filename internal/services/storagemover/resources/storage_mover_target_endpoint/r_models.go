// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_target_endpoint

type StorageMoverTargetEndpointModel struct {
	Name                 string `tfschema:"name"`
	StorageMoverId       string `tfschema:"storage_mover_id"`
	StorageAccountId     string `tfschema:"storage_account_id"`
	StorageContainerName string `tfschema:"storage_container_name"`
	Description          string `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &StorageMoverTargetEndpointModel{}
}
