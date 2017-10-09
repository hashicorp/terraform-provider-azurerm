package azurerm

import (
	"fmt"
	"strings"
)

func validateNetworkSecurityRuleProtocol(v interface{}, k string) (ws []string, errors []error) {
	value := strings.ToLower(v.(string))
	protocols := map[string]bool{
		"tcp": true,
		"udp": true,
		"*":   true,
	}

	if !protocols[value] {
		errors = append(errors, fmt.Errorf("Network Security Rule Protocol can only be Tcp, Udp or *"))
	}
	return
}
