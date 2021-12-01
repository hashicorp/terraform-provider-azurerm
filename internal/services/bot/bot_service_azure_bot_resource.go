package bot

import (
	"github.com/Azure/azure-sdk-for-go/services/botservice/mgmt/2021-03-01/botservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WebAppBotServiceResource struct {
	base botBaseResource
}

var _ sdk.Resource = WebAppBotServiceResource{}

var _ sdk.ResourceWithUpdate = WebAppBotServiceResource{}

func (r WebAppBotServiceResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
	return r.base.arguments(schema)
}

func (r WebAppBotServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r WebAppBotServiceResource) ModelObject() interface{} {
	return nil
}

func (r WebAppBotServiceResource) ResourceType() string{
	return "azurerm_bot_service_web_app"
}

func (r WebAppBotServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BotServiceID
}

func (r WebAppBotServiceResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), string(botservice.KindAzurebot))
}

func (r WebAppBotServiceResource) Read() sdk.ResourceFunc {
	return r.base.readFunc()
}

func (r WebAppBotServiceResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r WebAppBotServiceResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}
