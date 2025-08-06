// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package advisor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/advisor/2023-01-01/suppressions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AdvisorSuppressionResource struct{}

func TestAccAnalysisServicesServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_advisor_suppression", "test")
	r := AdvisorSuppressionResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (AdvisorSuppressionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := suppressions.ParseScopedSuppressionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Advisor.SuppressionsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil && resp.Model.Id != nil), nil
}

func (t AdvisorSuppressionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_advisor_recommendations" "test" {}

# The recommendation_names local variable is used to sort the recommendation names.
locals {
  recommendation_names = sort(data.azurerm_advisor_recommendations.test.recommendations[*].recommendation_name)
}

resource "azurerm_advisor_suppression" "test" {
  name              = "acctest%d"
  recommendation_id = local.recommendation_names[0]
  resource_id       = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ttl               = "00:30:00"
}
`, data.RandomInteger)
}
