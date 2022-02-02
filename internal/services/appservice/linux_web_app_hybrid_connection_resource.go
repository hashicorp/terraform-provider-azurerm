package appservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"strings"
)

type LinuxWebAppHybridConnectionResource struct{}

type LinuxWebAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = LinuxWebAppHybridConnectionResource{}

var _ sdk.ResourceWithCustomImporter = LinuxWebAppHybridConnectionResource{}

func (r LinuxWebAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r LinuxWebAppHybridConnectionResource) ResourceType() string {
	return "azurerm_linux_web_app_hybrid_connection"
}

func (r LinuxWebAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r LinuxWebAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r LinuxWebAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r LinuxWebAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r LinuxWebAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r LinuxWebAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r LinuxWebAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}

func (r LinuxWebAppHybridConnectionResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		os, sku, err := helpers.ServicePlanInfoForApp(ctx, metadata)
		if err != nil {
			return err
		}

		if os != nil && !strings.EqualFold(*os, "linux") {
			return fmt.Errorf("specified App is not a Linux App")
		}

		if helpers.PlanIsConsumption(*sku) || helpers.PlanIsElastic(*sku) {
			return fmt.Errorf("unsupported plan type. Hybrid Connections are not supported on Consumption or Elastic service plans")
		}

		return nil
	}
}
