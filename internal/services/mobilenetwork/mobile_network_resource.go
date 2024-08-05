// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mobilenetwork

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mobilenetwork/2022-11-01/mobilenetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MobileNetworkResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	MobileCountryCode string            `tfschema:"mobile_country_code"`
	MobileNetworkCode string            `tfschema:"mobile_network_code"`
	Tags              map[string]string `tfschema:"tags"`
	ServiceKey        string            `tfschema:"service_key"`
}

type MobileNetworkResource struct{}

var _ sdk.ResourceWithUpdate = MobileNetworkResource{}

func (r MobileNetworkResource) ResourceType() string {
	return "azurerm_mobile_network"
}

func (r MobileNetworkResource) ModelObject() interface{} {
	return &MobileNetworkResourceModel{}
}

func (r MobileNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return mobilenetwork.ValidateMobileNetworkID
}

func (r MobileNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"mobile_country_code": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{3}$`),
				"Mobile country code should be three digits.",
			),
		},

		"mobile_network_code": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{2,3}$`),
				"Mobile network code should be two or three digits.",
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (r MobileNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"service_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MobileNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MobileNetworkResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MobileNetwork.MobileNetworkClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := mobilenetwork.NewMobileNetworkID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := mobilenetwork.MobileNetwork{
				Location: location.Normalize(model.Location),
				Properties: mobilenetwork.MobileNetworkPropertiesFormat{
					PublicLandMobileNetworkIdentifier: mobilenetwork.PlmnId{
						Mcc: model.MobileCountryCode,
						Mnc: model.MobileNetworkCode,
					},
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MobileNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.MobileNetworkClient

			id, err := mobilenetwork.ParseMobileNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MobileNetworkResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			if metadata.ResourceData.HasChange("mobile_country_code") || metadata.ResourceData.HasChange("mobile_network_code") {
				model.Properties.PublicLandMobileNetworkIdentifier = mobilenetwork.PlmnId{
					Mcc: state.MobileCountryCode,
					Mnc: state.MobileNetworkCode,
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				model.Tags = &state.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MobileNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.MobileNetworkClient

			id, err := mobilenetwork.ParseMobileNetworkID(metadata.ResourceData.Id())
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

			state := MobileNetworkResourceModel{
				Name:              id.MobileNetworkName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.MobileCountryCode = model.Properties.PublicLandMobileNetworkIdentifier.Mcc
				state.MobileNetworkCode = model.Properties.PublicLandMobileNetworkIdentifier.Mnc

				if model.Properties.ServiceKey != nil {
					state.ServiceKey = *model.Properties.ServiceKey
				}

				if model.Tags != nil {
					state.Tags = *model.Tags
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MobileNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MobileNetwork.MobileNetworkClient

			id, err := mobilenetwork.ParseMobileNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// a workaround for that some child resources may still exist for seconds before it fully deleted.
			// tracked on https://github.com/Azure/azure-rest-api-specs/issues/22691
			// it will cause the error "Can not delete resource before nested resources are deleted."
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}
			stateConf := &pluginsdk.StateChangeConf{
				Delay:   5 * time.Minute,
				Pending: []string{"409"},
				Target:  []string{"200", "202"},
				Refresh: func() (result interface{}, state string, err error) {
					resp, err := client.Delete(ctx, *id)
					if err != nil {
						if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusConflict {
							return nil, "409", nil
						}
						return nil, "", err
					}
					return resp, "200", nil
				},
				MinTimeout: 15 * time.Second,
				Timeout:    time.Until(deadline),
			}

			if future, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for deleting of %s: %+v", id, err)
			} else {
				poller := future.(mobilenetwork.DeleteOperationResponse).Poller
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
