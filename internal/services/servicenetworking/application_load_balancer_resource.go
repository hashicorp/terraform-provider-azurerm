// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApplicationLoadBalancerResource struct{}

type ApplicationLoadBalancerModel struct {
	Name                          string            `tfschema:"name"`
	ResourceGroupName             string            `tfschema:"resource_group_name"`
	Location                      string            `tfschema:"location"`
	PrimaryConfigurationEndpoints string            `tfschema:"primary_configuration_endpoint"`
	Tags                          map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = ApplicationLoadBalancerResource{}

func (t ApplicationLoadBalancerResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9]$`), "the name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens. The value must be 1-64 characters long."),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (t ApplicationLoadBalancerResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"primary_configuration_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (t ApplicationLoadBalancerResource) ModelObject() interface{} {
	return &ApplicationLoadBalancerModel{}
}

func (t ApplicationLoadBalancerResource) ResourceType() string {
	return "azurerm_application_load_balancer"
}

func (t ApplicationLoadBalancerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trafficcontrollerinterface.ValidateTrafficControllerID
}

func (t ApplicationLoadBalancerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.TrafficControllerInterface
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ApplicationLoadBalancerModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			id := trafficcontrollerinterface.NewTrafficControllerID(subscriptionId, config.ResourceGroupName, config.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(t.ResourceType(), id)
			}

			payload := trafficcontrollerinterface.TrafficController{
				Location: location.Normalize(config.Location),
				Tags:     pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t ApplicationLoadBalancerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.TrafficControllerInterface

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", metadata.ResourceData.Id(), err)
			}

			state := ApplicationLoadBalancerModel{
				Name:              id.TrafficControllerName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					if endpoint := prop.ConfigurationEndpoints; endpoint != nil && len(*endpoint) > 0 {
						state.PrimaryConfigurationEndpoints = (*endpoint)[0]
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (t ApplicationLoadBalancerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.TrafficControllerInterface

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config ApplicationLoadBalancerModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			payload := trafficcontrollerinterface.TrafficControllerUpdate{
				Tags: pointer.To(config.Tags),
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (t ApplicationLoadBalancerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.TrafficControllerInterface

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// a workaround for that some child resources may still exist for seconds before it fully deleted.
			// tracked o https://github.com/Azure/azure-rest-api-specs/issues/26000
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
						if resp.HttpResponse.StatusCode == http.StatusConflict {
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
				poller := future.(trafficcontrollerinterface.DeleteOperationResponse).Poller
				if err := poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
