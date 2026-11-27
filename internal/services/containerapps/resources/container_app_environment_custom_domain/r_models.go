// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_custom_domain

type CustomDomainCertificateKeyVaultModel struct {
	Identity         string `tfschema:"identity"`
	KeyVaultSecretId string `tfschema:"key_vault_secret_id"`
}

type ContainerAppEnvironmentCustomDomainModel struct {
	ManagedEnvironmentId string `tfschema:"container_app_environment_id"`

	CertificatePassword string `tfschema:"certificate_password"`
	CertificateValue    string `tfschema:"certificate_blob_base64"`
	DnsSuffix           string `tfschema:"dns_suffix"`

	CertificateKeyVault []CustomDomainCertificateKeyVaultModel `tfschema:"certificate_key_vault"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppEnvironmentCustomDomainModel{}
}
