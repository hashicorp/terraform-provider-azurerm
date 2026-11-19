// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_app_hybrid_connection

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomImporter = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_web_app_hybrid_connection"
}

func (r Resource) Description() string {
	return "Manages a Web App Hybrid Connection."
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return webapps.ValidateRelayID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := webapps.ParseRelayID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		appId := commonids.NewAppServiceID(id.SubscriptionId, id.ResourceGroupName, id.SiteName)

		_, sku, err := helpers.ServicePlanInfoForApp(ctx, metadata, appId)
		if err != nil {
			return err
		}

		if helpers.PlanIsConsumption(sku) || helpers.PlanIsElastic(sku) {
			return fmt.Errorf("unsupported plan type. Hybrid Connections are not supported on Consumption or Elastic service plans")
		}

		return nil
	}
}
