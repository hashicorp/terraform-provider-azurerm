// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_storage

type ContainerAppEnvironmentStorageModel struct {
	Name                      string                                         `tfschema:"name"`
	ContainerAppEnvironmentId string                                         `tfschema:"container_app_environment_id"`
	AccountName               string                                         `tfschema:"account_name"`
	AccessKey                 string                                         `tfschema:"access_key"`
	AccessKeyVaultSecret      []ContainerAppEnvironmentStorageKeyVaultSecret `tfschema:"access_key_vault_secret"`
	ShareName                 string                                         `tfschema:"share_name"`
	AccessMode                string                                         `tfschema:"access_mode"`
	NfsServer                 string                                         `tfschema:"nfs_server_url"`
}

type ContainerAppEnvironmentStorageKeyVaultSecret struct {
	Identity         string `tfschema:"identity"`
	KeyVaultSecretId string `tfschema:"key_vault_secret_id"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentStorageModel{}
}
