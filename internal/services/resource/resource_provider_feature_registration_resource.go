// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name resource_provider_feature_registration -properties "name,provider_name" -known-values "subscription_id:data.Subscriptions.Secondary" -test-name "basicForResourceIdentity"

var _ sdk.ResourceWithIdentity = ResourceProviderFeatureRegistrationResource{}

type ResourceProviderFeatureRegistrationResource struct{}

func (r ResourceProviderFeatureRegistrationResource) Identity() resourceids.ResourceId {
	return new(features.FeatureId)
}

type ResourceProviderFeatureRegistrationModel struct {
	Name         string `tfschema:"name"`
	ProviderName string `tfschema:"provider_name"`
}

func (r ResourceProviderFeatureRegistrationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"provider_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceproviders.EnhancedValidate,
		},
	}
}

func (r ResourceProviderFeatureRegistrationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceProviderFeatureRegistrationResource) ModelObject() interface{} {
	return &ResourceProviderFeatureRegistrationModel{}
}

func (r ResourceProviderFeatureRegistrationResource) ResourceType() string {
	return "azurerm_resource_provider_feature_registration"
}

func (r ResourceProviderFeatureRegistrationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			var obj ResourceProviderFeatureRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			featureId := features.NewFeatureID(metadata.Client.Account.SubscriptionId, obj.ProviderName, obj.Name)

			provider, err := client.Get(ctx, featureId)
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("%s was not found", featureId)
				}

				return fmt.Errorf("retrieving %s: %+v", featureId, err)
			}

			registrationState := ""
			if model := provider.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
				registrationState = *model.Properties.State
			}

			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", featureId)
			}

			if strings.EqualFold(registrationState, Registered) {
				return metadata.ResourceRequiresImport(r.ResourceType(), featureId)
			}

			if strings.EqualFold(registrationState, Pending) {
				return fmt.Errorf("%s which requires manual approval can not be managed by Terraform", featureId)
			}

			resp, err := client.Register(ctx, featureId)
			if err != nil {
				return fmt.Errorf("registering %s: %+v", featureId, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
				if strings.EqualFold(*resp.Model.Properties.State, Pending) {
					return fmt.Errorf("%s which requires manual approval can not be managed by terraform", featureId)
				}
			}

			pollerType := custompollers.NewResourceProviderFeatureRegistrationPoller(client, featureId, Registered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for registration of %s: %+v", featureId, err)
			}

			metadata.SetID(featureId)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &featureId)
		},

		Timeout: 30 * time.Minute,
	}
}

func (r ResourceProviderFeatureRegistrationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			featureId, err := features.ParseFeatureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *featureId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(featureId)
				}

				return fmt.Errorf("retrieving %s: %+v", featureId, err)
			}

			registrationState := ""
			if model := resp.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
				registrationState = *model.Properties.State
			}

			if !strings.EqualFold(registrationState, Registered) {
				return metadata.MarkAsGone(featureId)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, featureId); err != nil {
				return err
			}

			return metadata.Encode(&ResourceProviderFeatureRegistrationModel{
				Name:         featureId.FeatureName,
				ProviderName: featureId.ProviderName,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ResourceProviderFeatureRegistrationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			featureId, err := features.ParseFeatureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Unregister(ctx, *featureId); err != nil {
				return fmt.Errorf("unregistering %s: %+v", featureId, err)
			}

			pollerType := custompollers.NewResourceProviderFeatureRegistrationPoller(client, *featureId, Unregistered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for registration of %s: %+v", featureId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceProviderFeatureRegistrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return features.ValidateFeatureID
}
