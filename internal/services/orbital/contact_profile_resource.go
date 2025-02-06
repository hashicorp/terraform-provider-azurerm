// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContactProfileResource struct{}

var _ sdk.ResourceWithUpdate = ContactProfileResource{}

var _ sdk.ResourceWithDeprecationAndNoReplacement = ContactProfileResource{}

func (r ContactProfileResource) DeprecationMessage() string {
	return "The `azurerm_orbital_contact_profile` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider."
}

type ContactProfileResourceModel struct {
	Name                           string                    `tfschema:"name"`
	ResourceGroup                  string                    `tfschema:"resource_group_name"`
	Location                       string                    `tfschema:"location"`
	MinimumVariableContactDuration string                    `tfschema:"minimum_variable_contact_duration"`
	MinimumElevationDegrees        float64                   `tfschema:"minimum_elevation_degrees"`
	AutoTrackingConfiguration      string                    `tfschema:"auto_tracking"`
	EventHubUri                    string                    `tfschema:"event_hub_uri"`
	Links                          []ContactProfileLinkModel `tfschema:"links"`
	Tags                           map[string]string         `tfschema:"tags"`
	NetworkConfigurationSubnetId   string                    `tfschema:"network_configuration_subnet_id"`
}

func (r ContactProfileResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"links": ContactProfileLinkSchema(),

		"network_configuration_subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"minimum_variable_contact_duration": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"minimum_elevation_degrees": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"auto_tracking": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(contactprofile.AutoTrackingConfigurationDisabled),
				string(contactprofile.AutoTrackingConfigurationSBand),
				string(contactprofile.AutoTrackingConfigurationXBand),
			}, false),
		},

		"event_hub_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": tags.Schema(),
	}
}

func (r ContactProfileResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContactProfileResource) ModelObject() interface{} {
	return &ContactProfileResourceModel{}
}

func (r ContactProfileResource) ResourceType() string {
	return "azurerm_orbital_contact_profile"
}

func (r ContactProfileResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ContactProfileResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Orbital.ContactProfileClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := contactprofile.NewContactProfileID(subscriptionId, model.ResourceGroup, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			links, err := expandContactProfileLinks(model.Links)
			if err != nil {
				return fmt.Errorf("expanding `links`: %+v", err)
			}

			autoTrackingConfiguration := contactprofile.AutoTrackingConfiguration(model.AutoTrackingConfiguration)

			networkConfiguration := contactprofile.ContactProfilesPropertiesNetworkConfiguration{
				SubnetId: model.NetworkConfigurationSubnetId,
			}

			// The service only accept `null` or non-empty value, empty string will cause a 400 response
			var eventHubUri *string
			if model.EventHubUri != "" {
				eventHubUri = pointer.To(model.EventHubUri)
			}

			contactProfilesProperties := contactprofile.ContactProfilesProperties{
				AutoTrackingConfiguration:    pointer.To(autoTrackingConfiguration),
				EventHubUri:                  eventHubUri,
				Links:                        links,
				MinimumElevationDegrees:      pointer.To(model.MinimumElevationDegrees),
				MinimumViableContactDuration: pointer.To(model.MinimumVariableContactDuration),
				NetworkConfiguration:         networkConfiguration,
			}

			contactProfile := contactprofile.ContactProfile{
				Id:         utils.String(id.ID()),
				Location:   model.Location,
				Name:       utils.String(model.Name),
				Properties: contactProfilesProperties,
				Tags:       pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, contactProfile); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContactProfileResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.ContactProfileClient
			id, err := contactprofile.ParseContactProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				state := ContactProfileResourceModel{
					Name:                           id.ContactProfileName,
					ResourceGroup:                  id.ResourceGroupName,
					Location:                       model.Location,
					MinimumVariableContactDuration: pointer.From(props.MinimumViableContactDuration),
					MinimumElevationDegrees:        pointer.From(props.MinimumElevationDegrees),
					AutoTrackingConfiguration:      string(pointer.From(props.AutoTrackingConfiguration)),
					EventHubUri:                    pointer.From(props.EventHubUri),
					NetworkConfigurationSubnetId:   props.NetworkConfiguration.SubnetId,
				}
				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
				links, err := flattenContactProfileLinks(props.Links)
				if err != nil {
					return err
				}
				state.Links = links

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r ContactProfileResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.ContactProfileClient
			id, err := contactprofile.ParseContactProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContactProfileResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return contactprofile.ValidateContactProfileID
}

func (r ContactProfileResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.ContactProfileClient
			id, err := contactprofile.ParseContactProfileID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContactProfileResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			contactProfileLinks, err := expandContactProfileLinks(state.Links)
			if err != nil {
				return err
			}

			autoTrackingConfiguration := contactprofile.AutoTrackingConfiguration(state.AutoTrackingConfiguration)

			networkConfiguration := contactprofile.ContactProfilesPropertiesNetworkConfiguration{
				SubnetId: state.NetworkConfigurationSubnetId,
			}

			// The service only accept `null` or non-empty value, empty string will cause a 400 response
			var eventHubUri *string
			if state.EventHubUri != "" {
				eventHubUri = pointer.To(state.EventHubUri)
			}

			if metadata.ResourceData.HasChangesExcept("name", "resource_group_name") {
				contactProfile := contactprofile.ContactProfile{
					Location: state.Location,
					Properties: contactprofile.ContactProfilesProperties{
						AutoTrackingConfiguration:    pointer.To(autoTrackingConfiguration),
						EventHubUri:                  eventHubUri,
						Links:                        contactProfileLinks,
						MinimumElevationDegrees:      pointer.To(state.MinimumElevationDegrees),
						MinimumViableContactDuration: pointer.To(state.MinimumVariableContactDuration),
						NetworkConfiguration:         networkConfiguration,
					},
					Tags: pointer.To(state.Tags),
				}

				if err := client.CreateOrUpdateThenPoll(ctx, *id, contactProfile); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}
