package appservice

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WebAppActiveSlotResource ActiveSlotResource

var _ sdk.Resource = WebAppActiveSlotResource{}

func (r WebAppActiveSlotResource) ModelObject() interface{} {
	return &ActiveSlotModel{}
}

func (r WebAppActiveSlotResource) ResourceType() string {
	return "azurerm_web_app_active_slot"
}

func (r WebAppActiveSlotResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.WebAppSlotID
}

func (r WebAppActiveSlotResource) Arguments() map[string]*pluginsdk.Schema {
	return ActiveSlotResource{}.Arguments()
}

func (r WebAppActiveSlotResource) Attributes() map[string]*pluginsdk.Schema {
	return ActiveSlotResource{}.Attributes()
}

func (r WebAppActiveSlotResource) Create() sdk.ResourceFunc {
	return ActiveSlotResource{}.Create()
}

func (r WebAppActiveSlotResource) Read() sdk.ResourceFunc {
	return ActiveSlotResource{}.Read()
}

func (r WebAppActiveSlotResource) Delete() sdk.ResourceFunc {
	return ActiveSlotResource{}.Delete()
}
