package eventgrid

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = EventGridPartnerRegistrationResource{}

type EventGridPartnerRegistrationResource struct{}

type EventGridPartnerRegistrationResourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	PartnerRegistrationID string            `tfschema:"partner_registration_id"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (EventGridPartnerRegistrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotEmpty,
				validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"EventGrid partner registration name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			),
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"tags": commonschema.Tags(),
	}
}

func (EventGridPartnerRegistrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"partner_registration_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (EventGridPartnerRegistrationResource) ModelObject() interface{} {
	return &EventGridPartnerRegistrationResourceModel{}
}

func (EventGridPartnerRegistrationResource) ResourceType() string {
	return "azurerm_eventgrid_partner_registration"
}

func (r EventGridPartnerRegistrationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EventGridPartnerRegistrationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := partnerregistrations.NewPartnerRegistrationID(subscriptionId, config.ResourceGroup, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := partnerregistrations.PartnerRegistration{
				Location: "global",
				Name:     pointer.To(config.Name),
				Tags:     pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r EventGridPartnerRegistrationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations

			id, err := partnerregistrations.ParsePartnerRegistrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerRegistrationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r EventGridPartnerRegistrationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations

			id, err := partnerregistrations.ParsePartnerRegistrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := EventGridPartnerRegistrationResourceModel{
				Name:          id.PartnerRegistrationName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.PartnerRegistrationID = pointer.From(props.PartnerRegistrationImmutableId)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r EventGridPartnerRegistrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations

			id, err := partnerregistrations.ParsePartnerRegistrationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (EventGridPartnerRegistrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return partnerregistrations.ValidatePartnerRegistrationID
}
