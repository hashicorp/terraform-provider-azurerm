// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func IsManagedHSMURI(env environments.Environment, uri string) (bool, error, string, string) {
	url, err := url.Parse(uri)
	if err != nil {
		return false, fmt.Errorf("Error parsing %s as URI: %+v", uri, err), "", ""
	}

	instanceName, domainSuffix, found := strings.Cut(url.Hostname(), ".")
	if !found {
		return false, fmt.Errorf("Key vault URI hostname does not have the right number of components: %s", url.Hostname()), "", ""
	}
	expectedDomainSuffix, found := env.ManagedHSM.DomainSuffix()
	if !found {
		expectedDomainSuffix = utils.String("managedhsm.azure.net")
	}
	if domainSuffix == *expectedDomainSuffix {
		return true, nil, instanceName, domainSuffix
	} else {
		return false, nil, "", ""
	}
}
