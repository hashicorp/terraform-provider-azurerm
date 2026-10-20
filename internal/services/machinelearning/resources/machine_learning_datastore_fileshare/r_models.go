// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_fileshare

type MachineLearningDataStoreFileShareModel struct {
	Name                  string            `tfschema:"name"`
	WorkSpaceID           string            `tfschema:"workspace_id"`
	StorageFileShareID    string            `tfschema:"storage_fileshare_id"`
	Description           string            `tfschema:"description"`
	IsDefault             bool              `tfschema:"is_default"`
	ServiceDataIdentity   string            `tfschema:"service_data_identity"`
	AccountKey            string            `tfschema:"account_key"`
	SharedAccessSignature string            `tfschema:"shared_access_signature"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MachineLearningDataStoreFileShareModel{}
}
