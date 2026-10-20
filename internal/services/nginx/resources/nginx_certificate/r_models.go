// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_certificate

type CertificateModel struct {
	Name                   string `tfschema:"name"`
	NginxDeploymentId      string `tfschema:"nginx_deployment_id"`
	KeyVirtualPath         string `tfschema:"key_virtual_path"`
	CertificateVirtualPath string `tfschema:"certificate_virtual_path"`
	KeyVaultSecretId       string `tfschema:"key_vault_secret_id"`
}

func (m Resource) ModelObject() any {
	return &CertificateModel{}
}
