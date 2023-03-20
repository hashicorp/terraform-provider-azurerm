package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

func TestAccOrchestratedVMSSDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_orchestrated_virtual_machine_scale_set", "test")
	d := OrchestratedVirtualMachineScaleSetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("network_interfaces.#").HasValue("1"),
			),
		},
	})
}

func (OrchestratedVirtualMachineScaleSetDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data azurerm_orchestrated_virtual_machine_scale_set test {
  name                = azurerm_service_plan.test.name
  resource_group_name = azurerm_service_plan.test.resource_group_name
}
`, OrchestratedVirtualMachineScaleSetResource{}.linuxInstances(data))
}
