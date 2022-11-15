package bot

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/botservice/mgmt/2021-05-01-preview/botservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AzureBotServiceResource struct {
	base botBaseResource
}

var _ sdk.ResourceWithUpdate = AzureBotServiceResource{}

var _ sdk.ResourceWithCustomImporter = AzureBotServiceResource{}

func (r AzureBotServiceResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r AzureBotServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r AzureBotServiceResource) ModelObject() interface{} {
	return nil
}

func (r AzureBotServiceResource) ResourceType() string {
	return "azurerm_bot_service_azure_bot"
}

func (r AzureBotServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BotServiceID
}

func (r AzureBotServiceResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), string(botservice.KindAzurebot))
}

func (r AzureBotServiceResource) Read() sdk.ResourceFunc {
	return r.base.readFunc()
}

func (r AzureBotServiceResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r AzureBotServiceResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}

func (r AzureBotServiceResource) CustomImporter() sdk.ResourceRunFunc {
	return r.base.importerFunc(string(botservice.KindAzurebot))
}
