// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_blobstorage

type MachineLearningDataStoreBlobStorageModel struct {
	Name                    string            `tfschema:"name"`
	WorkSpaceID             string            `tfschema:"workspace_id"`
	StorageContainerID      string            `tfschema:"storage_container_id"`
	Description             string            `tfschema:"description"`
	IsDefault               bool              `tfschema:"is_default"`
	ServiceDataAuthIdentity string            `tfschema:"service_data_auth_identity"`
	AccountKey              string            `tfschema:"account_key"`
	SharedAccessSignature   string            `tfschema:"shared_access_signature"`
	Tags                    map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MachineLearningDataStoreBlobStorageModel{}
}
