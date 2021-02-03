package eventhub

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/sdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ConsumerGroupObject struct {
	Name              string `tfschema:"name"`
	NamespaceName     string `tfschema:"namespace_name"`
	EventHubName      string `tfschema:"eventhub_name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	UserMetadata      string `tfschema:"user_metadata"`
}

var _ sdk.Resource = ConsumerGroupResource{}
var _ sdk.ResourceWithUpdate = ConsumerGroupResource{}

type ConsumerGroupResource struct {
}

func (r ConsumerGroupResource) ResourceType() string {
	return "azurerm_eventhub_consumer_group"
}

func (r ConsumerGroupResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateEventHubConsumerName(),
		},

		"namespace_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateEventHubNamespaceName(),
		},

		"eventhub_name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: ValidateEventHubName(),
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"user_metadata": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 1024),
		},
	}
}

func (r ConsumerGroupResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r ConsumerGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state ConsumerGroupObject
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("creating Consumer Group %q..", state.Name)
			client := metadata.Client.Eventhub.ConsumerGroupClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := parse.NewEventHubConsumerGroupID(subscriptionId, state.ResourceGroupName, state.NamespaceName, state.EventHubName, state.Name)
			existing, err := client.Get(ctx, state.ResourceGroupName, state.NamespaceName, state.EventHubName, state.Name)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of an existing Consumer Group %q: %+v", state.Name, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := eventhub.ConsumerGroup{
				Name: utils.String(state.Name),
				ConsumerGroupProperties: &eventhub.ConsumerGroupProperties{
					UserMetadata: utils.String(state.UserMetadata),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, state.ResourceGroupName, state.NamespaceName, state.EventHubName, state.Name, parameters); err != nil {
				return fmt.Errorf("creating Consumer Group %q (EventHub %q / Namespace %q / Resource Group %q): %+v", state.Name, state.EventHubName, state.NamespaceName, state.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ConsumerGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.EventHubConsumerGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Info("Decoding state..")
			var state ConsumerGroupObject
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("updating Consumer Group %q..", state.Name)
			client := metadata.Client.Eventhub.ConsumerGroupClient

			parameters := eventhub.ConsumerGroup{
				Name: utils.String(id.ConsumergroupName),
				ConsumerGroupProperties: &eventhub.ConsumerGroupProperties{
					UserMetadata: utils.String(state.UserMetadata),
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.ConsumergroupName, parameters); err != nil {
				return fmt.Errorf("updating Consumer Group %q (EventHub %q / Namespace %q / Resource Group %q): %+v", id.ConsumergroupName, id.EventhubName, id.NamespaceName, id.ResourceGroup, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ConsumerGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Eventhub.ConsumerGroupClient
			id, err := parse.EventHubConsumerGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving Consumer Group %q..", id.ConsumergroupName)
			resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.ConsumergroupName)
			if err != nil {
				return fmt.Errorf("reading Consumer Group %q (EventHub %q / Namespace %q / Resource Group %q): %+v", id.ConsumergroupName, id.EventhubName, id.NamespaceName, id.ResourceGroup, err)
			}

			state := ConsumerGroupObject{
				Name:              id.ConsumergroupName,
				NamespaceName:     id.NamespaceName,
				EventHubName:      id.EventhubName,
				ResourceGroupName: id.ResourceGroup,
			}

			if props := resp.ConsumerGroupProperties; props != nil {
				state.UserMetadata = utils.NormalizeNilableString(props.UserMetadata)
			}

			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ConsumerGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Eventhub.ConsumerGroupClient
			id, err := parse.EventHubConsumerGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting Consumer Group %q..", id.ConsumergroupName)
			if resp, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.EventhubName, id.ConsumergroupName); err != nil {
				if !utils.ResponseWasNotFound(resp) {
					return fmt.Errorf("deleting Consumer Group %q (EventHub %q / Namespace %q / Resource Group %q): %+v", id.ConsumergroupName, id.EventhubName, id.NamespaceName, id.ResourceGroup, err)
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ConsumerGroupResource) ModelObject() interface{} {
	return ConsumerGroupObject{}
}

func (r ConsumerGroupResource) IDValidationFunc() schema.SchemaValidateFunc {
	return validate.EventHubConsumerGroupID
}
