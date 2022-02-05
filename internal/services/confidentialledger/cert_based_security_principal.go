package confidentialledger

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
)

func expandConfidentialLedgerCertBasedSecurityPrincipal(input []interface{}) *[]confidentialledger.CertBasedSecurityPrincipal {
	output := make([]confidentialledger.CertBasedSecurityPrincipal, len(input))
	if len(input) == 0 {
		return &output
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		cert := v["cert"].(string)
		ledgerRoleName := v["ledger_role_name"].(confidentialledger.LedgerRoleName)

		result := confidentialledger.CertBasedSecurityPrincipal{
			Cert:           &cert,
			LedgerRoleName: &ledgerRoleName,
		}

		output = append(output, result)
	}

	return &output
}

func flattenConfidentialLedgerCertBasedSecurityPrincipal(input *[]confidentialledger.CertBasedSecurityPrincipal) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)

	for _, item := range *input {
		cert := ""
		ledgerRoleName := confidentialledger.LedgerRoleNameReader

		if item.Cert != nil {
			cert = *item.Cert
		}

		if item.LedgerRoleName != nil {
			ledgerRoleName = *item.LedgerRoleName
		}

		aadBasedSecurityPrincipal := map[string]interface{}{
			"cert":             cert,
			"ledger_role_name": ledgerRoleName,
		}

		output = append(output, aadBasedSecurityPrincipal)
	}

	return &output, nil
}
