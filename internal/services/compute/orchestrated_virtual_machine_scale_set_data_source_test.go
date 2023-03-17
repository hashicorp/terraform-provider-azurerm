package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

func TestAccOrchestratedVMSSDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_service_plan", "test")
	d := OrchestratedVirtualMachineScaleSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
			),
		},
	})
}

func (OrchestratedVirtualMachineScaleSetDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_service_plan test {
  name                = azurerm_service_plan.test.name
  resource_group_name = azurerm_service_plan.test.resource_group_name
}
`, OrchestratedVirtualMachineScaleSetResource{}.linuxInstances(data))
}
