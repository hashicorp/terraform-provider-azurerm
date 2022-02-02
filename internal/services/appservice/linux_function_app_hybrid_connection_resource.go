package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LinuxFunctionAppHybridConnectionResource struct{}

type LinuxFunctionAppHybridConnectionModel struct {
	*AppHybridConnectionCommonModel
}

var _ sdk.ResourceWithUpdate = LinuxFunctionAppHybridConnectionResource{}

func (r LinuxFunctionAppHybridConnectionResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r LinuxFunctionAppHybridConnectionResource) ResourceType() string {
	return "azurerm_linux_function_app_hybrid_connection"
}

func (r LinuxFunctionAppHybridConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return AppHybridConnectionCommonResource{}.IDValidationFunc()
}

func (r LinuxFunctionAppHybridConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Arguments()
}

func (r LinuxFunctionAppHybridConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return AppHybridConnectionCommonResource{}.Attributes()
}

func (r LinuxFunctionAppHybridConnectionResource) Create() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Create()
}

func (r LinuxFunctionAppHybridConnectionResource) Read() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Read()
}

func (r LinuxFunctionAppHybridConnectionResource) Delete() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Delete()
}

func (r LinuxFunctionAppHybridConnectionResource) Update() sdk.ResourceFunc {
	return AppHybridConnectionCommonResource{}.Update()
}
