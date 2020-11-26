package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
)

func MsSqlDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.MsSqlDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Database resource id: %v", k, err))
	}

	return warnings, errors
}

func MsSqlDatabaseAutoPauseDelay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
		return warnings, errors
	}
	min := 60
	max := 10080
	if (v < min || v > max) && v%10 != 0 && v != -1 {
		errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d) and divisible by 10 or -1, got %d", k, min, max, v))
		return warnings, errors
	}

	return warnings, errors
}

func MsSqlDBSkuName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`(?i)(^(GP_S_Gen5_(1|2|4|6|8|10|12|14|16|18|20|24|32|40))$|^((GP|HS|BC)_Gen4_(1|2|3|4|5|6|7|8|9|10|16|24))$|^((GP|HS|BC)_Gen5_(2|4|6|8|10|12|14|16|18|20|24|32|40|80))$|^(BC_M_(8|10|12|14|16|18|20|24|32|64|128))$|^(Basic)$|^(ElasticPool)$|^(S(0|1|2|3|4|6|7|9|12))$|^(P(1|2|4|6|11|15))$|^(DW(1|2|3|4|5|6|7|8|9)000*c)$|^(DS(1|2|3|4|5|6|10|12|15|20)00)$)`),

		`This is not a valid sku name. For example, a valid sku name is 'GP_S_Gen5_1','HS_Gen4_1','BC_Gen5_2', 'ElasticPool', 'Basic', 'S0', 'P1'.`,
	)
}

func MsSqlDBCollation() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`(^[A-Z]+)([A-Za-z0-9]+_)+((BIN|BIN2|CI_AI|CI_AI_KS|CI_AI_KS_WS|CI_AI_WS|CI_AS|CI_AS_KS|CI_AS_KS_WS|CS_AI|CS_AI_KS|CS_AI_KS_WS|CS_AI_WS|CS_AS|CS_AS_KS|CS_AS_KS_WS|CS_AS_WS)+)((_[A-Za-z0-9]+)+$)*`),

		`This is not a valid collation.`,
	)
}

func MsSqlRestorableDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.RestorableDroppedDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Restorable Database resource id: %v", k, err))
	}

	return warnings, errors
}

func MsSqlRecoverableDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.RecoverableDBID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Recoverable Database resource id: %v", k, err))
	}

	return warnings, errors
}
