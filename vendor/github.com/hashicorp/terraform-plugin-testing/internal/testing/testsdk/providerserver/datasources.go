// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/datasource"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/provider"
)

func ProviderDataSource(p provider.Provider, typeName string) (datasource.DataSource, *tfprotov6.Diagnostic) {
	d, ok := p.DataSourcesMap()[typeName]

	if !ok {
		return nil, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Missing Data Source Type",
			Detail:   "The provider does not define the data source type: " + typeName,
		}
	}

	return d, nil
}
