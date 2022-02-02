package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WindowsFunctionAppHybridConnectionResource struct{}

type WindowsFunctionAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = WindowsFunctionAppHybridConnectionResource{}

func (r WindowsFunctionAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r WindowsFunctionAppHybridConnectionResource) ResourceType() string {
	return "azurerm_windows_function_app_hybrid_connection"
}

func (r WindowsFunctionAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r WindowsFunctionAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r WindowsFunctionAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r WindowsFunctionAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r WindowsFunctionAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r WindowsFunctionAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r WindowsFunctionAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}
