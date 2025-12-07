package mssql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

type MsSqlExecuteJobAction struct{}

func TestAccMsSqlExecuteJobAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_execute_job", "test")
	a := MsSqlExecuteJobAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.basic(data),
				Check:  nil, // TODO
			},
		},
	})
}

func TestAccMsSqlExecuteJobAction_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_execute_job", "test")
	a := MsSqlExecuteJobAction{}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		Steps: []resource.TestStep{
			{
				Config: a.complete(data),
				Check:  nil, // TODO
			},
		},
	})
}

func (r *MsSqlExecuteJobAction) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

action "azurerm_mssql_execute_job" "test" {
  config {
    job_id = azurerm_mssql_job.test.id
  }
}
`, r.template(data))
}

func (r *MsSqlExecuteJobAction) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

action "azurerm_mssql_execute_job" "test" {
  config {
    job_id              = azurerm_mssql_job.test.id
    wait_for_completion = true
    timeout             = "5m"
  }
}
`, r.template(data))
}

func (r *MsSqlExecuteJobAction) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "terraform_data" "trigger" {
  input = azurerm_mssql_job_step.test.id
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_mssql_execute_job.test]
    }
  }
}
`, MsSqlJobStepTestResource{}.basic(data))
}
