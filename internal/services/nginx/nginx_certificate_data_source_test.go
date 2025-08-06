// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nginx_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type NginxCertificateDataSource struct{}

func TestAccNginxCertificateDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_nginx_certificate", "test")
	r := NginxCertificateDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("certificate_virtual_path").Exists(),
				check.That(data.ResourceName).Key("key_vault_secret_id").Exists(),
				check.That(data.ResourceName).Key("key_virtual_path").Exists(),
				check.That(data.ResourceName).Key("sha1_thumbprint").Exists(),
				check.That(data.ResourceName).Key("key_vault_secret_version").Exists(),
				check.That(data.ResourceName).Key("key_vault_secret_creation_date").Exists(),
			),
		},
	})
}

func (d NginxCertificateDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_nginx_certificate" "test" {
  name                = azurerm_nginx_certificate.test.name
  nginx_deployment_id = azurerm_nginx_deployment.test.id
}
`, CertificateResource{}.basic(data))
}
