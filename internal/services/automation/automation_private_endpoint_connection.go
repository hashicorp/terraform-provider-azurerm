package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateEndpointConnectionModel struct {
	ResourceGroupName     string   `tfschema:"resource_group_name"`
	AutomationAccountName string   `tfschema:"automation_account_name"`
	Name                  string   `tfschema:"name"`
	LinkStatus            string   `tfschema:"link_status"`
	LinkDescription       string   `tfschema:"link_description"`
	LinkActionRequired    string   `tfschema:"link_action_required"`
	PrivateEndpointID     string   `tfschema:"private_endpoint_id"`
	GroupIds              []string `tfschema:"group_ids"`
}

type PrivateEndpointConnectionResource struct{}

var _ sdk.Resource = (*PrivateEndpointConnectionResource)(nil)

func (m PrivateEndpointConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),
		"automation_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"link_status": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"link_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (m PrivateEndpointConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"group_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"link_action_required": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"private_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m PrivateEndpointConnectionResource) ModelObject() interface{} {
	return &PrivateEndpointConnectionModel{}
}

func (m PrivateEndpointConnectionResource) ResourceType() string {
	return "azurerm_automation_private_endpoint_connection"
}

func (m PrivateEndpointConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.PrivateEndpointClient

			var model PrivateEndpointConnectionModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := parse.NewPrivateEndpointConnectionID(subscriptionID, model.ResourceGroupName,
				model.AutomationAccountName, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return fmt.Errorf("retreiving %s: %v", id, err)
			}
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("no such private endpoint connection: %v", id)
			}

			var params automation.PrivateEndpointConnection
			params.PrivateEndpointConnectionProperties = &automation.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: &automation.PrivateLinkServiceConnectionStateProperty{},
			}
			if model.LinkStatus != "" {
				params.PrivateLinkServiceConnectionState.Status = utils.String(model.LinkStatus)
			}
			if model.LinkDescription != "" {
				params.PrivateLinkServiceConnectionState.Description = utils.String(model.LinkDescription)
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, params)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}
			// TODO may not need wait, delete lines below
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m PrivateEndpointConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.PrivateEndpointConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.PrivateEndpointClient
			result, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name)
			if err != nil {
				return err
			}

			var output PrivateEndpointConnectionModel
			output.Name = id.Name
			output.ResourceGroupName = id.ResourceGroup
			output.AutomationAccountName = id.AutomationAccountName
			if result.PrivateEndpoint != nil {
				output.PrivateEndpointID = utils.NormalizeNilableString(result.PrivateEndpoint.ID)
			}
			if prop := result.PrivateEndpointConnectionProperties; prop != nil {
				if state := prop.PrivateLinkServiceConnectionState; state != nil {
					output.LinkStatus = utils.NormalizeNilableString(state.Status)
					output.LinkDescription = utils.NormalizeNilableString(state.Description)
					output.LinkActionRequired = utils.NormalizeNilableString(state.ActionsRequired)
				}
				if point := prop.PrivateEndpoint; point != nil {
					output.PrivateEndpointID = utils.NormalizeNilableString(point.ID)
				}
				// no group id in definition
			}
			return meta.Encode(&output)
		},
	}
}

func (m PrivateEndpointConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.PrivateEndpointClient

			id, err := parse.PrivateEndpointConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateEndpointConnectionModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd automation.PrivateEndpointConnection
			upd.PrivateEndpointConnectionProperties = &automation.PrivateEndpointConnectionProperties{}
			upd.PrivateEndpointConnectionProperties.PrivateLinkServiceConnectionState = &automation.PrivateLinkServiceConnectionStateProperty{}
			if meta.ResourceData.HasChange("link_status") {
				upd.PrivateLinkServiceConnectionState.Status = utils.String(model.LinkStatus)
			}
			if meta.ResourceData.HasChange("link_description") {
				upd.PrivateLinkServiceConnectionState.Description = utils.String(model.LinkDescription)
			}
			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m PrivateEndpointConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.PrivateEndpointConnectionID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.PrivateEndpointClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m PrivateEndpointConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.PrivateEndpointConnectionID
}
