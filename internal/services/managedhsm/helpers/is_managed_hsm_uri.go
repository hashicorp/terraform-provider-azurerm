// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func IsManagedHSMURI(env environments.Environment, uri string) (bool, string, string, error) {
	// TODO: remove this function once https://github.com/hashicorp/terraform-provider-azurerm/pull/26521
	// is merged (`./internal/customermanagedkeys/key_vault_or_managed_hsm_key.go` replaces this)
	expectedDomainSuffix, ok := env.ManagedHSM.DomainSuffix()
	if !ok {
		return false, "", "", errors.New("managed HSM isn't supported in this environment")
	}

	url, err := url.Parse(uri)
	if err != nil {
		return false, "", "", fmt.Errorf("parsing %s as URI: %+v", uri, err)
	}

	instanceName, domainSuffix, found := strings.Cut(url.Hostname(), ".")
	if !found {
		return false, "", "", fmt.Errorf("key vault URI hostname does not have the right number of components: %s", url.Hostname())
	}
	if domainSuffix == *expectedDomainSuffix {
		return true, instanceName, domainSuffix, nil
	} else {
		return false, "", "", nil
	}
}
