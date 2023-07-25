// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type HealthCareFhirServiceDataSource struct{}

func TestAccHealthCareFhirServiceDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
	r := HealthCareFhirServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func TestAccHealthCareFhirServiceDataSource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_fhir_service", "test")
	r := HealthCareFhirServiceDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.dataSourceUserAssignedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists()),
		},
	})
}

func (HealthCareFhirServiceDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_fhir_service" "test" {
  name         = azurerm_healthcare_fhir_service.test.name
  workspace_id = azurerm_healthcare_fhir_service.test.workspace_id
}
`, HealthcareApiFhirServiceResource{}.basic(data))
}

func (HealthCareFhirServiceDataSource) dataSourceUserAssignedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_healthcare_fhir_service" "test" {
  name         = azurerm_healthcare_fhir_service.test.name
  workspace_id = azurerm_healthcare_fhir_service.test.workspace_id
}
`, HealthcareApiFhirServiceResource{}.updateIdentityUserAssigned(data))
}
