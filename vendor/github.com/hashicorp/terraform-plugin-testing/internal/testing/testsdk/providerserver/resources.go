// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
)

func ProviderResource(p provider.Provider, typeName string) (resource.Resource, *tfprotov6.Diagnostic) {
	r, ok := p.ResourcesMap()[typeName]

	if !ok {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Missing Resource Type",
			Detail:   "The provider does not define the resource type: " + typeName,
		}
	}

	return r, nil
}
