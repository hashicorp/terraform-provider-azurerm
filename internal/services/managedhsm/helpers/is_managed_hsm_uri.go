// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func IsManagedHSMURI(env environments.Environment, uri string) (bool, error, string, string) {
	// TODO: remove this function once https://github.com/hashicorp/terraform-provider-azurerm/pull/26521
	// is merged (`./internal/customermanagedkeys/key_vault_or_managed_hsm_key.go` replaces this)
	expectedDomainSuffix, ok := env.ManagedHSM.DomainSuffix()
	if !ok {
		return false, fmt.Errorf("Managed HSM isn't sipported in this environment"), "", ""
	}

	url, err := url.Parse(uri)
	if err != nil {
		return false, fmt.Errorf("Error parsing %s as URI: %+v", uri, err), "", ""
	}

	instanceName, domainSuffix, found := strings.Cut(url.Hostname(), ".")
	if !found {
		return false, fmt.Errorf("Key vault URI hostname does not have the right number of components: %s", url.Hostname()), "", ""
	}
	if domainSuffix == *expectedDomainSuffix {
		return true, nil, instanceName, domainSuffix
	} else {
		return false, nil, "", ""
	}
}
