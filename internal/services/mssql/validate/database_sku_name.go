// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	Basic            = "(Basic)"
	Elastic          = "(ElasticPool)"
	Standard         = "(S(0|1|2|3|4|6|7|9|12))"
	Premium          = "(P(1|2|4|6|11|15))"
	DataWarehouse    = "(DW(1|2|3|4|5|6|7|8|9)5?000*c)"
	Stretch          = "(DS(1|2|3|4|5|6|10|12|15|20)00)"
	BusinessCritical = "(BC_M_(8|10|12|14|16|18|20|24|32|64|128))"
	Gen4             = "((GP|HS|BC)_Gen4_(1|2|3|4|5|6|7|8|9|10|16|24))"
	Gen5             = "(GP|HS|BC)_Gen5_(2|4|6|8|10|12|14|16|18|20|24|32|40|80|128)"
	ServerlessGen5   = "(GP|HS)_S_Gen5_(1|2|4|6|8|10|12|14|16|18|20|24|32|40|80)"
	Fsv2             = "(GP_Fsv2_(8|10|12|14|16|18|20|24|32|36|72))"
	Dc               = "((GP|BC|HS)_DC_(2|4|6|8|10|12|14|16|18|20|32|40))"
	EightIM          = "(HS_8IM_(24|48|80))"
	Serverless8IM    = "(HS_S_8IM_(24|80))"
	Premium8IM       = "(HS_PRMS_(2|4|6|8|10|12|14|16|18|20|24|32|40|64|80|128))"
	Premium8IH       = "(HS_MOPRMS_(2|4|6|8|10|12|14|16|18|20|24|32|40|64|80))"
)

func DatabaseSkuName() pluginsdk.SchemaValidateFunc {
	pattern := "(?i)(^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$)"
	return validation.StringMatch(
		regexp.MustCompile(fmt.Sprintf(pattern, Basic, Elastic, Standard, Premium, DataWarehouse, Stretch, BusinessCritical, Gen4, Gen5, ServerlessGen5, Fsv2, Dc, EightIM, Serverless8IM, Premium8IM, Premium8IH)),

		`This is not a valid sku name. For example, a valid sku name is 'GP_S_Gen5_1','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'S0', 'P1'.`,
	)
}

// DatabaseSkuNameWithFree accepts the retired `Free` SKU with a deprecation warning for v5.x compatibility.
// TODO 6.0: remove once `features.SixPointOh()` is removed.
func DatabaseSkuNameWithFree() pluginsdk.SchemaValidateFunc {
	baseValidate := DatabaseSkuName()
	return func(i any, k string) ([]string, []error) {
		if v, ok := i.(string); ok && strings.EqualFold(v, "Free") {
			return []string{"the `Free` SKU has been retired by Azure and will be removed in v6.0 of the AzureRM Provider - use a serverless General Purpose SKU (for example `GP_S_Gen5_2`) with `free_limit_enabled = true` instead"}, nil
		}

		return baseValidate(i, k)
	}
}
