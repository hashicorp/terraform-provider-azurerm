package costmanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2022-06-01-preview/scheduledactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AnomalyAlertResource struct{}

func TestAccResourceAnomalyAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_costmanagement_anomaly_alert", "test")
	testResource := AnomalyAlertResource{}
	data.ResourceTest(t, testResource, []acceptance.TestStep{
		data.ApplyStep(testResource.basicConfig, testResource),
		data.ImportStep(),
	})
}

// go install && make acctests SERVICE='costmanagement' TESTARGS='-run=TestAccResourceAnomalyAlert_basic'

func (AnomalyAlertResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scheduledactions.ParseScopedScheduledActionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.CostManagement.ScheduledActionsClient.GetByScope(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model.Properties != nil), nil
}

func (AnomalyAlertResource) basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_costmanagement_anomaly_alert" "test" {
  name            = "acctestRG-%d"
  email_subject   = "Hi"
  email_addresses = ["test@test.com", "test@hashicorp.developer"]
	message         = "Oops, cost anomaly"
}
`, data.RandomInteger)
}
