// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/statestore"
)

func ProviderStateStore(p provider.Provider, typeName string) (statestore.StateStore, *tfprotov6.Diagnostic) {
	s, ok := p.StateStoresMap()[typeName]

	if !ok {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Missing State Store Type",
			Detail:   "The provider does not define the state store type: " + typeName,
		}
	}

	return s, nil
}
