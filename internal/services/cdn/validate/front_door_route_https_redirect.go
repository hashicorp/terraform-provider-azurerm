// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func SupportsBothHttpAndHttps(input []interface{}, key string) error {
	if len(input) == 0 {
		return fmt.Errorf("expected %q to be a list of string", key)
	}

	for _, str := range input {
		_, ok := str.(string)
		if !ok {
			return fmt.Errorf("expected %q to contain only strings", key)
		}
	}

	protocols := utils.ExpandStringSlice(input)
	if !utils.SliceContainsValue(*protocols, string(endpoints.DestinationProtocolHTTP)) || !utils.SliceContainsValue(*protocols, string(endpoints.DestinationProtocolHTTPS)) {
		return fmt.Errorf("'https_redirect_enabled' and 'supported_protocols' conflict. The 'https_redirect_enabled' field cannot be set to 'true' unless the 'supported_protocols' field contains both 'Http' and 'Https'")
	}

	return nil
}
