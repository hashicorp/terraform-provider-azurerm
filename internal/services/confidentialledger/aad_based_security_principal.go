package confidentialledger

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandConfidentialLedgerAADBasedSecurityPrincipal(input []interface{}) *[]confidentialledger.AADBasedSecurityPrincipal {
	output := make([]confidentialledger.AADBasedSecurityPrincipal, 0, len(input))
	if len(input) == 0 {
		return &output
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		ledgerRoleName := confidentialledger.LedgerRoleName(v["ledger_role_name"].(string))
		principalId := v["principal_id"].(string)
		tenantId := v["tenant_id"].(string)

		result := confidentialledger.AADBasedSecurityPrincipal{
			LedgerRoleName: &ledgerRoleName,
			PrincipalId:    utils.String(principalId),
			TenantId:       utils.String(tenantId),
		}

		output = append(output, result)
	}

	return &output
}

func flattenConfidentialLedgerAADBasedSecurityPrincipal(input *[]confidentialledger.AADBasedSecurityPrincipal) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0, len(*input))

	for _, item := range *input {
		ledgerRoleName := string(confidentialledger.LedgerRoleNameReader)
		principalId := ""
		tenantId := ""

		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		}

		if item.PrincipalId != nil {
			principalId = *item.PrincipalId
		}

		if item.TenantId != nil {
			tenantId = *item.TenantId
		}

		aadBasedSecurityPrincipal := map[string]interface{}{
			"ledger_role_name": ledgerRoleName,
			"principal_id":     principalId,
			"tenant_id":        tenantId,
		}

		output = append(output, aadBasedSecurityPrincipal)
	}

	return &output, nil
}
