package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2022-02-01/availabilitygrouplisteners"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MsSqlVirtualMachineAvailabilityGroupListenerResource struct{}

func TestAccMsSqlVirtualMachineAvailabilityGroupListenerResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_virtual_machine_availability_group_listener", "test")
	r := MsSqlVirtualMachineAvailabilityGroupListenerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (MsSqlVirtualMachineAvailabilityGroupListenerResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {

	id, err := availabilitygrouplisteners.ParseAvailabilityGroupListenerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.VirtualMachinesAvailabilityGroupListenersClient.Get(ctx, *id, availabilitygrouplisteners.GetOperationOptions{Expand: utils.String("*")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("%s does not exist", *id)
		}
		return nil, fmt.Errorf("reading %s: %v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

// TODO add proper test
func (r MsSqlVirtualMachineAvailabilityGroupListenerResource) basic() string {
	return fmt.Sprintf(`
resource "azurerm_mssql_virtual_machine_availability_group_listener" "test" {
  name                              = "MyListener"
  resource_group_name               = // *** resource_group_name
  availability_group_name           = "default2"
  port                              = 1432
  create_default_availability_group = true
  sql_virtual_machine_group_name    = // *** sql_virtual_machine_group_name


  load_balancer_configuration {
    private_ip_address {
      ip_address = "10.0.2.8"
      subnet_id  = 
    }

    load_balancer_id              = // *** load_balancer_id
    probe_port                    = 51572
    sql_virtual_machine_instances = [ // *** sql_virtual_machine_instances ]
  }

  replica {
    sql_virtual_machine_instance_id = // *** sql_virtual_machine_instance_id
    role                            = "PRIMARY"
    commit                          = "SYNCHRONOUS_COMMIT"
    failover                        = "AUTOMATIC"
    readable_secondary              = "ALL"
  }
}
`)
}
