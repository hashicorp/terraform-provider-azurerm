package chaosstudio

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ChaosStudioTargetResource{}

type ChaosStudioTargetResource struct{}

type ChaosStudioTargetResourceSchema struct {
	Location         string `tfschema:"location"`
	TargetType       string `tfschema:"target_type"`
	TargetResourceId string `tfschema:"target_resource_id"`
}

func (r ChaosStudioTargetResource) ModelObject() interface{} {
	return nil
}

func (r ChaosStudioTargetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return nil
}

func (r ChaosStudioTargetResource) ResourceType() string {
	return ""
}

func (r ChaosStudioTargetResource) Arguments() map[string]*pluginsdk.Schema {
	return nil
}

func (r ChaosStudioTargetResource) Attributes() map[string]*pluginsdk.Schema {
	return nil
}

func (r ChaosStudioTargetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}

func (r ChaosStudioTargetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{}
}
