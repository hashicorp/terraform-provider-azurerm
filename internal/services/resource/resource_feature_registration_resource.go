// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ResourceProviderFeatureRegistrationResource{}

type ResourceProviderFeatureRegistrationResource struct{}

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
	return "azurerm_resource_feature_registration"
}

func (r ResourceProviderFeatureRegistrationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient
			account := metadata.Client.Account

			var obj ResourceProviderFeatureRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			featureId := features.NewFeatureID(account.SubscriptionId, obj.ProviderName, obj.Name)

			provider, err := client.Get(ctx, featureId)
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("%s was not found", featureId)
				}

				return fmt.Errorf("retrieving %q: %+v", featureId, err)
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
				return fmt.Errorf("%s which requires manual approval can not be managed by terraform", featureId)
			}

			log.Printf("[DEBUG] Registering %s..", featureId)
			resp, err := client.Register(ctx, featureId)
			if err != nil {
				return fmt.Errorf("error registering feature %q: %+v", featureId, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
				if strings.EqualFold(*resp.Model.Properties.State, Pending) {
					return fmt.Errorf("%s which requires manual approval can not be managed by terraform", featureId)
				}
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{Registering},
				Target:     []string{Registered},
				Refresh:    r.featureRegisteringStateRefreshFunc(ctx, client, featureId),
				MinTimeout: 3 * time.Minute,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s registering to be completed: %+v", featureId, err)
			}

			log.Printf("[DEBUG] Registered Resource Provider %q.", featureId)

			metadata.SetID(featureId)
			return nil
		},

		Timeout: 120 * time.Minute,
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
				log.Printf("[WARN] %s was not registered - removing from state", *featureId)
				return metadata.MarkAsGone(featureId)
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
				return fmt.Errorf("error checking for existing feature %q: %+v", *featureId, err)
			}

			existing, err := client.Get(ctx, *featureId)
			if err != nil {
				return fmt.Errorf("get feature %q: %+v", *featureId, err)
			}

			if existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.State != nil {
				if strings.EqualFold(*existing.Model.Properties.State, Pending) {
					return fmt.Errorf("%s which requires manual approval should not be managed by terraform", *featureId)
				}
				if strings.EqualFold(*existing.Model.Properties.State, Unregistered) {
					return nil
				}
			}

			log.Printf("[INFO] unregistering feature %q.", *featureId)

			resp, err := client.Unregister(ctx, *featureId)
			if err != nil {
				return fmt.Errorf("unregistering feature %q: %+v", *featureId, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
				if strings.EqualFold(*resp.Model.Properties.State, Pending) {
					return fmt.Errorf("%s requires manual registration approval and can not be managed by terraform", *featureId)
				}
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context had no deadline")
			}
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{Unregistering},
				Target:     []string{NotRegistered, Unregistered},
				Refresh:    r.featureRegisteringStateRefreshFunc(ctx, client, *featureId),
				MinTimeout: 3 * time.Minute,
				Timeout:    time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be complete unregistering: %+v", *featureId, err)
			}

			log.Printf("[DEBUG] Unregistered Feature %q.", *featureId)

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r ResourceProviderFeatureRegistrationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return features.ValidateFeatureID
}

func (r ResourceProviderFeatureRegistrationResource) featureRegisteringStateRefreshFunc(ctx context.Context, client *features.FeaturesClient, id features.FeatureId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.State == nil {
			return nil, "", fmt.Errorf("error reading %s registering status: %+v", id, err)
		}

		return res, *res.Model.Properties.State, nil
	}
}
