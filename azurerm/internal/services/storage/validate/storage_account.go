package validate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
)

func StorageAccountID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.StorageAccountID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func StorageAccountIPRule(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|30))?$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must start with IPV4 address and/or slash, number of bits (0-30) as prefix. Example: 23.45.1.0/30.", k))
		return warnings, errors
	}

	ipParts := strings.Split(v.(string), ".")
	firstIPPart := ipParts[0]
	secondIPPart, _ := strconv.Atoi(ipParts[1])
	if (firstIPPart == "10") || (firstIPPart == "172" && secondIPPart >= 16 && secondIPPart <= 31) || (firstIPPart == "192" && secondIPPart == 168) {
		errors = append(errors, fmt.Errorf("%q must be public ip address", k))
		return warnings, errors
	}

	return warnings, errors
}
