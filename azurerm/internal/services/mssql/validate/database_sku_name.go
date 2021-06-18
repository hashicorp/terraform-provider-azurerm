package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func DatabaseSkuName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`(?i)(^(GP_S_Gen5_(1|2|4|6|8|10|12|14|16|18|20|24|32|40))$|^((GP|HS|BC)_Gen4_(1|2|3|4|5|6|7|8|9|10|16|24))$|^((GP|HS|BC)_Gen5_(2|4|6|8|10|12|14|16|18|20|24|32|40|80))$|^(BC_M_(8|10|12|14|16|18|20|24|32|64|128))$|^(Basic)$|^(ElasticPool)$|^(S(0|1|2|3|4|6|7|9|12))$|^(P(1|2|4|6|11|15))$|^(DW(1|2|3|4|5|6|7|8|9)5?000*c)$|^(DS(1|2|3|4|5|6|10|12|15|20)00)$)`),

		`This is not a valid sku name. For example, a valid sku name is 'GP_S_Gen5_1','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'S0', 'P1'.`,
	)
}
