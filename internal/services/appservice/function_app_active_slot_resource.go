package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FunctionAppActiveSlotResource ActiveSlotResource

var _ sdk.Resource = FunctionAppActiveSlotResource{}

func (r FunctionAppActiveSlotResource) ModelObject() interface{} {
	return &ActiveSlotModel{}
}

func (r FunctionAppActiveSlotResource) ResourceType() string {
	return "azurerm_function_app_active_slot"
}

func (r FunctionAppActiveSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppSlotID
}

func (r FunctionAppActiveSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return ActiveSlotResource{}.Arguments()
}

func (r FunctionAppActiveSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return ActiveSlotResource{}.Attributes()
}

func (r FunctionAppActiveSlotResource) Create() sdk.ResourceFunc {
	return ActiveSlotResource{}.Create()
}

func (r FunctionAppActiveSlotResource) Read() sdk.ResourceFunc {
	return ActiveSlotResource{}.Read()
}

func (r FunctionAppActiveSlotResource) Delete() sdk.ResourceFunc {
	return ActiveSlotResource{}.Delete()
}
