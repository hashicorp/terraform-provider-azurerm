package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsWebAppHybridConnectionResource struct{}

type WindowsWebAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = WindowsWebAppHybridConnectionResource{}

func (r WindowsWebAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r WindowsWebAppHybridConnectionResource) ResourceType() string {
	return "azurerm_windows_web_app_hybrid_connection"
}

func (r WindowsWebAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r WindowsWebAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r WindowsWebAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r WindowsWebAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r WindowsWebAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r WindowsWebAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r WindowsWebAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}
