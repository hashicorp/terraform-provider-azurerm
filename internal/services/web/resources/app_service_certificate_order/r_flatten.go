// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_certificate_order

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/certificateregistration/2023-12-01/appservicecertificateorders"
)

func flattenArmCertificateOrderCertificate(input *map[string]appservicecertificateorders.AppServiceCertificate) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for k, v := range *input {
		result := make(map[string]any)

		result["certificate_name"] = k

		if keyVaultID := v.KeyVaultId; keyVaultID != nil {
			result["key_vault_id"] = *keyVaultID
		}
		if keyVaultSecretName := v.KeyVaultSecretName; keyVaultSecretName != nil {
			result["key_vault_secret_name"] = *keyVaultSecretName
		}
		result["provisioning_state"] = pointer.FromEnum(v.ProvisioningState)

		results = append(results, result)
	}

	return results
}

func flattenAppServiceCertificateNotRenewableReasons(input *[]appservicecertificateorders.ResourceNotRenewableReason) []string {
	results := make([]string, 0)

	if input == nil {
		return results
	}

	for _, v := range *input {
		results = append(results, string(v))
	}

	return results
}

func flattenProductType(input appservicecertificateorders.CertificateProductType) string {
	if input == appservicecertificateorders.CertificateProductTypeStandardDomainValidatedWildCardSsl {
		return "WildCard"
	}
	return "Standard"
}
