// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/list"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
)

func ProviderListResource(p provider.Provider, typeName string) (list.ListResource, *tfprotov6.Diagnostic) {
	r, ok := p.ListResourcesMap()[typeName]

	if !ok {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Missing List Resource Type",
			Detail:   "The provider does not define the list resource type: " + typeName,
		}
	}

	return r, nil
}
