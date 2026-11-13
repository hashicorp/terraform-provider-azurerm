// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
)

func ExpandApiManagementCommonHostnameConfiguration(input map[string]any, hostnameType apimanagementservice.HostnameType) apimanagementservice.HostnameConfiguration {
	output := apimanagementservice.HostnameConfiguration{
		Type: hostnameType,
	}
	if v, ok := input["certificate"]; ok && v.(string) != "" {
		output.EncodedCertificate = pointer.To(v.(string))
	}
	if v, ok := input["certificate_password"]; ok && v.(string) != "" {
		output.CertificatePassword = pointer.To(v.(string))
	}
	if v, ok := input["host_name"]; ok && v.(string) != "" {
		output.HostName = v.(string)
	}
	if v, ok := input["key_vault_certificate_id"]; ok && v.(string) != "" {
		output.KeyVaultId = pointer.To(v.(string))
	}

	if v, ok := input["negotiate_client_certificate"]; ok {
		output.NegotiateClientCertificate = pointer.To(v.(bool))
	}

	if v, ok := input["ssl_keyvault_identity_client_id"]; ok && v.(string) != "" {
		output.IdentityClientId = pointer.To(v.(string))
	}

	return output
}
