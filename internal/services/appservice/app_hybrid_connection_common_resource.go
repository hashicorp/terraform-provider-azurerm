package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/2017-04-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/2017-04-01/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AppHybridConnectionCommonResource struct{}

type AppHybridConnectionCommonModel struct {
	AppServiceId        string `tfschema:"app_id"`
	RelayId             string `tfschema:"relay_id"`
	HostName            string `tfschema:"hostname"`
	HostPort            int    `tfschema:"port"`
	SendKeyName         string `tfschema:"send_key_name"`
	NamespaceName       string `tfschema:"namespace_name"`
	RelayName           string `tfschema:"relay_name"`
	ServiceBusNamespace string `tfschema:"service_bus_namespace"`
	ServiceBusSuffix    string `tfschema:"service_bus_suffix"`
	SendKeyValue        string `tfschema:"send_key_value"`
}

var _ sdk.ResourceWithUpdate = AppHybridConnectionCommonResource{}

func (r AppHybridConnectionCommonResource) ModelObject() interface{} {
	return &AppHybridConnectionCommonModel{}
}

func (r AppHybridConnectionCommonResource) ResourceType() string {
	return "" // never called
}

func (r AppHybridConnectionCommonResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AppHybridConnectionID
}

func (r AppHybridConnectionCommonResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WebAppID,
		},

		"relay_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hybridconnections.ValidateHybridConnectionID,
		},

		"hostname": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: azValidate.PortNumberOrZero,
		},

		"send_key_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "RootManageSharedAccessKey",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r AppHybridConnectionCommonResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"relay_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_bus_namespace": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_bus_suffix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"send_key_value": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r AppHybridConnectionCommonResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var appHybridConn AppHybridConnectionCommonModel

			client := metadata.Client.AppService.WebAppsClient

			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}
			appId, err := parse.WebAppID(appHybridConn.AppServiceId)
			if err != nil {
				return err
			}
			relayId, err := hybridconnections.ParseHybridConnectionID(appHybridConn.RelayId)
			if err != nil {
				return err
			}

			id := parse.NewAppHybridConnectionID(appId.SubscriptionId, appId.ResourceGroup, appId.SiteName, relayId.NamespaceName, relayId.HybridConnectionName)

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if !utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}
			if existing.ID != nil && *existing.ID != "" {
				return tf.ImportAsExistsError(r.ResourceType(), *existing.ID)
			}

			envelope := web.HybridConnection{
				HybridConnectionProperties: &web.HybridConnectionProperties{
					RelayArmURI: utils.String(relayId.ID()),
					Hostname:    utils.String(appHybridConn.HostName),
					Port:        utils.Int32(int32(appHybridConn.HostPort)),
					SendKeyName: utils.String(appHybridConn.SendKeyName),
				},
			}
			key, err := getSendKeyValue(ctx, metadata, id, appHybridConn.SendKeyName)
			envelope.HybridConnectionProperties.SendKeyValue = key

			_, err = client.CreateOrUpdateHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName, envelope)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r AppHybridConnectionCommonResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			appHybridConn := AppHybridConnectionCommonModel{
				AppServiceId:  parse.NewWebAppID(id.SubscriptionId, id.ResourceGroup, id.SiteName).ID(),
				RelayName:     id.RelayName,
				NamespaceName: id.HybridConnectionNamespaceName,
			}

			if props := existing.HybridConnectionProperties; props != nil {
				appHybridConn.RelayId = utils.NormalizeNilableString(props.RelayArmURI)
				appHybridConn.HostName = utils.NormalizeNilableString(props.Hostname)
				appHybridConn.HostPort = int(utils.NormaliseNilableInt32(props.Port))
				appHybridConn.SendKeyName = utils.NormalizeNilableString(existing.SendKeyName)
				appHybridConn.ServiceBusNamespace = utils.NormalizeNilableString(props.ServiceBusNamespace)
				appHybridConn.ServiceBusSuffix = utils.NormalizeNilableString(props.ServiceBusSuffix)
				appHybridConn.SendKeyValue = utils.NormalizeNilableString(props.SendKeyValue)
			}

			if appHybridConn.ServiceBusNamespace != "" && appHybridConn.SendKeyName != "" {
				relayClient := metadata.Client.Relay.NamespacesClient
				connectionId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroup, id.HybridConnectionNamespaceName, appHybridConn.SendKeyName)
				keys, err := relayClient.ListKeys(ctx, connectionId)
				if err != nil {
					return err
				}
				if keys.Model != nil {
					appHybridConn.SendKeyValue = utils.NormalizeNilableString(keys.Model.PrimaryKey)
				}
			}

			return metadata.Encode(&appHybridConn)
		},
	}
}

func (r AppHybridConnectionCommonResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.DeleteHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if !response.WasNotFound(resp.Response) {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (r AppHybridConnectionCommonResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.WebAppsClient

			id, err := parse.AppHybridConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var appHybridConn AppHybridConnectionCommonModel
			if err := metadata.Decode(&appHybridConn); err != nil {
				return err
			}

			existing, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("hostname") {
				existing.HybridConnectionProperties.Hostname = utils.String(appHybridConn.HostName)
			}

			if metadata.ResourceData.HasChange("port") {
				existing.HybridConnectionProperties.Port = utils.Int32(int32(appHybridConn.HostPort))
			}

			if metadata.ResourceData.HasChange("send_key_name") {
				key, err := getSendKeyValue(ctx, metadata, *id, appHybridConn.SendKeyName)
				if err != nil {
					return err
				}
				existing.HybridConnectionProperties.SendKeyValue = key
			}

			_, err = client.CreateOrUpdateHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func getSendKeyValue(ctx context.Context, metadata sdk.ResourceMetaData, id parse.AppHybridConnectionId, sendKeyName string) (*string, error) {
	relayClient := metadata.Client.Relay.NamespacesClient
	connectionId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, id.ResourceGroup, id.HybridConnectionNamespaceName, sendKeyName)
	keys, err := relayClient.ListKeys(ctx, connectionId)
	if err != nil {
		return nil, err
	}
	if keys.Model == nil || keys.Model.PrimaryKey == nil {
		return nil, fmt.Errorf("reading Send Key Value for %s in %s", connectionId.AuthorizationRuleName, id)
	}
	return keys.Model.PrimaryKey, nil
}
