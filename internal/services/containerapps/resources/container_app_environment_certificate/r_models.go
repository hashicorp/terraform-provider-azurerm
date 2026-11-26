// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

type CertificateKeyVaultModel struct {
	Identity         string `tfschema:"identity"`
	KeyVaultSecretId string `tfschema:"key_vault_secret_id"`
}

type ContainerAppCertificateModel struct {
	Name                 string         `tfschema:"name"`
	ManagedEnvironmentId string         `tfschema:"container_app_environment_id"`
	Tags                 map[string]any `tfschema:"tags"`

	// Write only?
	CertificatePassword string `tfschema:"certificate_password"`
	CertificateBlob     string `tfschema:"certificate_blob_base64"`

	CertificateKeyVault []CertificateKeyVaultModel `tfschema:"certificate_key_vault"`

	// Read Only
	SubjectName    string `tfschema:"subject_name"`
	Issuer         string `tfschema:"issuer"`
	IssueDate      string `tfschema:"issue_date"`
	ExpirationDate string `tfschema:"expiration_date"`
	Thumbprint     string `tfschema:"thumbprint"`
}

func (r Resource) ModelObject() any {
	return &ContainerAppCertificateModel{}
}
