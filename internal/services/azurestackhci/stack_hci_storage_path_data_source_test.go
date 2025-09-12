// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type StackHCIStoragePathDataSource struct{}

func TestAccStackHCIStoragePathDataSource_basic(t *testing.T) {
	if os.Getenv(customLocationIdEnv) == "" {
		t.Skipf("skipping since %q has not been specified", customLocationIdEnv)
	}

	data := acceptance.BuildTestData(t, "data.azurerm_stack_hci_storage_path", "test")
	d := StackHCIStoragePathDataSource{}

	data.DataSourceTestInSequence(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").IsNotEmpty(),
				check.That(data.ResourceName).Key("custom_location_id").IsNotEmpty(),
				check.That(data.ResourceName).Key("path").IsNotEmpty(),
			),
		},
	})
}

func (d StackHCIStoragePathDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_stack_hci_storage_path" "test" {
  name                = azurerm_stack_hci_storage_path.test.name
  resource_group_name = azurerm_stack_hci_storage_path.test.resource_group_name
}
`, StackHCIStoragePathResource{}.complete(data))
}
