package confidentialledger

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
)

func expandConfidentialLedgerCertBasedSecurityPrincipal(input []interface{}) *[]confidentialledger.CertBasedSecurityPrincipal {
	output := make([]confidentialledger.CertBasedSecurityPrincipal, 0, len(input))
	if len(input) == 0 {
		return &output
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		cert := v["cert"].(string)
		ledgerRoleName := confidentialledger.LedgerRoleName(v["ledger_role_name"].(string))

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

	output := make([]interface{}, 0, len(*input))

	for _, item := range *input {
		var cert string
		var ledgerRoleName string

		if item.Cert != nil {
			cert = *item.Cert
		} else {
			return nil, fmt.Errorf("error flattening cert-based security principal (no certificate): %v", item)
		}

		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		} else {
			return nil, fmt.Errorf("error flattening cert-based security principal (no assigned role): %v", item)
		}

		aadBasedSecurityPrincipal := map[string]interface{}{
			"cert":             cert,
			"ledger_role_name": ledgerRoleName,
		}

		output = append(output, aadBasedSecurityPrincipal)
	}

	return &output, nil
}
