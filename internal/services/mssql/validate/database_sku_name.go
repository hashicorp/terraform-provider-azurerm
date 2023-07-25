// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	Free             = "(Free)"
	Basic            = "(Basic)"
	Elastic          = "(ElasticPool)"
	Standard         = "(S(0|1|2|3|4|6|7|9|12))"
	Premium          = "(P(1|2|4|6|11|15))"
	DataWarehouse    = "(DW(1|2|3|4|5|6|7|8|9)5?000*c)"
	Stretch          = "(DS(1|2|3|4|5|6|10|12|15|20)00)"
	BusinessCritical = "(BC_M_(8|10|12|14|16|18|20|24|32|64|128))"
	Gen4             = "((GP|HS|BC)_Gen4_(1|2|3|4|5|6|7|8|9|10|16|24))"
	Gen5             = "(GP|HS|BC)_Gen5_(2|4|6|8|10|12|14|16|18|20|24|32|40|80)"
	ServerlessGen5   = "(GP|HS)_S_Gen5_(1|2|4|6|8|10|12|14|16|18|20|24|32|40|80)"
	Fsv2             = "(GP_Fsv2_(8|10|12|14|16|18|20|24|32|36|72))"
	Dc               = "((GP|BC|HS)_DC_(2|4|6|8))"
	EightIM          = "(HS_8IM_(24|48|80))"
	Serverless8IM    = "(HS_S_8IM_(24|80))"
)

func DatabaseSkuName() pluginsdk.SchemaValidateFunc {
	pattern := "(?i)(^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$|^%s$)"
	return validation.StringMatch(regexp.MustCompile(fmt.Sprintf(pattern, Free, Basic, Elastic, Standard, Premium, DataWarehouse, Stretch, BusinessCritical, Gen4, Gen5, ServerlessGen5, Fsv2, Dc, EightIM, Serverless8IM)),

		`This is not a valid sku name. For example, a valid sku name is 'GP_S_Gen5_1','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'S0', 'P1'.`,
	)
}
