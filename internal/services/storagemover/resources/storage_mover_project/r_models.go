// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_project

type StorageMoverProjectResourceModel struct {
	Name           string `tfschema:"name"`
	StorageMoverId string `tfschema:"storage_mover_id"`
	Description    string `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &StorageMoverProjectResourceModel{}
}
