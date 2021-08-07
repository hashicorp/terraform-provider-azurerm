package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	free             = "(Free)"
	basic            = "(Basic)"
	elastic          = "(ElasticPool)"
	standard         = "(S(0|1|2|3|4|6|7|9|12))"
	premium          = "(P(1|2|4|6|11|15))"
	dataWarehouse    = "(DW(1|2|3|4|5|6|7|8|9)5?000*c)"
	stretch          = "(DS(1|2|3|4|5|6|10|12|15|20)00)"
	businessCritical = "(BC_M_(8|10|12|14|16|18|20|24|32|64|128))"
	gen4             = "((GP|HS|BC)_Gen4_(1|2|3|4|5|6|7|8|9|10|16|24))"
	gen5             = "(GP|HS|BC)_Gen5_(2|4|6|8|10|12|14|16|18|20|24|32|40|80)"
	serverlessGen5   = "(GP_S_Gen5_(1|2|4|6|8|10|12|14|16|18|20|24|32|40))"
	fsv2             = "(GP_Fsv2_(8|10|12|14|16|18|20|24|32|36|72))"
)

func DatabaseSkuName() pluginsdk.SchemaValidateFunc {
	pattern := "(?i)(^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$)"
	return validation.StringMatch(regexp.MustCompile(fmt.Sprintf(pattern, free, basic, elastic, standard, premium, dataWarehouse, stretch, businessCritical, gen4, gen5, serverlessGen5, fsv2)),

		`This is not a valid sku name. For example, a valid sku name is 'GP_S_Gen5_1','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'S0', 'P1'.`,
	)
}
