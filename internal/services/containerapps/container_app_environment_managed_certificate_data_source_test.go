// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentManagedCertificateDataSource struct{}

func TestAccContainerAppEnvironmentManagedCertificateDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment_managed_certificate", "test")
	r := ContainerAppEnvironmentManagedCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subject_name").IsSet(),
				check.That(data.ResourceName).Key("domain_control_validation_type").HasValue("CNAME"),
			),
		},
	})
}

func (d ContainerAppEnvironmentManagedCertificateDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment_managed_certificate" "test" {
  name                         = azurerm_container_app_environment_managed_certificate.test.name
  container_app_environment_id = azurerm_container_app_environment_managed_certificate.test.container_app_environment_id
}


`, ContainerAppEnvironmentManagedCertificateResource{}.basic(data))
}
