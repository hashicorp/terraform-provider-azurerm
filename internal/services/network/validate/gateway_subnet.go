// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

func IsGatewaySubnet(i interface{}, k string) (warnings []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := commonids.ParseSubnetIDInsensitively(value)
	if err != nil {
		errors = append(errors, err)
		return
	}

	if !strings.EqualFold(id.SubnetName, "GatewaySubnet") {
		errors = append(errors, fmt.Errorf("expected %s to reference a gateway subnet with name GatewaySubnet", k))
	}

	return warnings, errors
}
