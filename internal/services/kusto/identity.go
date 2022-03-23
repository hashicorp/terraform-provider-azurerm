package kusto

import (
	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2021-08-27/kusto"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandTrustedExternalTenants(input []interface{}) *[]kusto.TrustedExternalTenant {
	output := make([]kusto.TrustedExternalTenant, 0)

	for _, v := range input {
		if !features.ThreePointOhBeta() {
			if v.(string) == "MyTenantOnly" {
				return &[]kusto.TrustedExternalTenant{}
			}
		}
		output = append(output, kusto.TrustedExternalTenant{
			Value: utils.String(v.(string)),
		})
	}

	return &output
}

func flattenTrustedExternalTenants(input *[]kusto.TrustedExternalTenant) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	if !features.ThreePointOhBeta() && len(*input) == 0 {
		return append(output, "MyTenantOnly")
	}

	for _, v := range *input {
		if v.Value == nil {
			continue
		}

		output = append(output, *v.Value)
	}

	return output
}
