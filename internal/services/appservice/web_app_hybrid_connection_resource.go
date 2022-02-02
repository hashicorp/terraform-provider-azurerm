package appservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WebAppHybridConnectionResource struct{}

type WebAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = WebAppHybridConnectionResource{}

var _ sdk.ResourceWithCustomImporter = WebAppHybridConnectionResource{}

func (r WebAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r WebAppHybridConnectionResource) ResourceType() string {
	return "azurerm_web_app_hybrid_connection"
}

func (r WebAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r WebAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r WebAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r WebAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r WebAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r WebAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r WebAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}

func (r WebAppHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		_, sku, err := helpers.ServicePlanInfoForApp(ctx, metadata)
		if err != nil {
			return err
		}

		if helpers.PlanIsConsumption(*sku) || helpers.PlanIsElastic(*sku) {
			return fmt.Errorf("unsupported plan type. Hybrid Connections are not supported on Consumption or Elastic service plans")
		}

		return nil
	}
}
