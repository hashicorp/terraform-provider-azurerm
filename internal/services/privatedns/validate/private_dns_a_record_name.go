package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func PrivateDnsARecordName(i interface{}, k string) (warnings []string, errors []error) {
	validate.LowerCasedString(i, k)

	if strings.ContainsAny(i.(string), "@") {
		return nil, []error{fmt.Errorf("%q cannot contain @", k)}
	}

	return nil, nil
}
