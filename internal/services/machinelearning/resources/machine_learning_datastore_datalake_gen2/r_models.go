// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_datastore_datalake_gen2

type MachineLearningDataStoreDataLakeGen2Model struct {
	Name                string            `tfschema:"name"`
	WorkSpaceID         string            `tfschema:"workspace_id"`
	StorageContainerID  string            `tfschema:"storage_container_id"`
	TenantID            string            `tfschema:"tenant_id"`
	ClientID            string            `tfschema:"client_id"`
	ClientSecret        string            `tfschema:"client_secret"`
	AuthorityURL        string            `tfschema:"authority_url"`
	Description         string            `tfschema:"description"`
	IsDefault           bool              `tfschema:"is_default"`
	ServiceDataIdentity string            `tfschema:"service_data_identity"`
	Tags                map[string]string `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &MachineLearningDataStoreDataLakeGen2Model{}
}
