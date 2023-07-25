// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type ContainerAppEnvironmentCertificateDataSource struct{}

func TestAccContainerAppEnvironmentCertificateDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_container_app_environment_certificate", "test")
	r := ContainerAppEnvironmentCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("subject_name").IsSet(),
				check.That(data.ResourceName).Key("thumbprint").IsSet(),
				check.That(data.ResourceName).Key("issue_date").IsSet(),
				check.That(data.ResourceName).Key("expiration_date").IsSet(),
				check.That(data.ResourceName).Key("issuer").IsSet(),
			),
		},
	})
}

func (d ContainerAppEnvironmentCertificateDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_container_app_environment_certificate" "test" {
  name                         = azurerm_container_app_environment_certificate.test.name
  container_app_environment_id = azurerm_container_app_environment_certificate.test.container_app_environment_id
}


`, ContainerAppEnvironmentCertificateResource{}.basic(data))
}
