package appservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FunctionAppHybridConnectionResource struct{}

type FunctionAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = FunctionAppHybridConnectionResource{}

var _ sdk.ResourceWithCustomImporter = FunctionAppHybridConnectionResource{}

func (r FunctionAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r FunctionAppHybridConnectionResource) ResourceType() string {
	return "azurerm_function_app_hybrid_connection"
}

func (r FunctionAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r FunctionAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r FunctionAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r FunctionAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r FunctionAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r FunctionAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r FunctionAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}

func (r FunctionAppHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
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
