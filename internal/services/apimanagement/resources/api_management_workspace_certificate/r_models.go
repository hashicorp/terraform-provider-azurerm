// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_certificate

type ApiManagementWorkspaceCertificateModel struct {
	Name                         string `tfschema:"name"`
	ApiManagementWorkspaceId     string `tfschema:"api_management_workspace_id"`
	CertificateDataBase64        string `tfschema:"certificate_data_base64"`
	KeyVaultSecretId             string `tfschema:"key_vault_secret_id"`
	Password                     string `tfschema:"password"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`

	// Computed fields
	Expiration string `tfschema:"expiration"`
	Subject    string `tfschema:"subject"`
	Thumbprint string `tfschema:"thumbprint"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementWorkspaceCertificateModel{}
}
